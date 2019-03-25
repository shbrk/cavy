package net

import (
	"core/log"
	"errors"
	"net"
	"sync/atomic"
	"time"
)

func NewTCPConnection(conn *net.TCPConn,
	session ISession,
	writeQueueSize int,
	readBufferSize int,
	writeBufferSize int, readTimeout int) *TCPConnection {

	if writeBufferSize < MAX_PACKET_LENGTH {
		writeBufferSize = MAX_PACKET_LENGTH * 2
	}
	if readBufferSize < MAX_PACKET_LENGTH {
		readBufferSize = MAX_PACKET_LENGTH * 2
	}
	if writeQueueSize <= 0 {
		writeQueueSize = 1024 * 8
	}
	if readTimeout <= 0 {
		readTimeout = 10
	}

	var c = &TCPConnection{conn: conn,
		state:       k_UNINITED,
		session:     session,
		writeQueue:  make(chan *Packet, writeQueueSize),
		readBuffer:  NewByteBuffer(readBufferSize),
		writeBuffer: NewByteBuffer(writeBufferSize),
		readTimeout: time.Duration(readTimeout) * time.Second}
	_ = conn.SetNoDelay(true)
	c.init()
	return c
}

const (
	k_UNINITED    = 0
	k_ESTABLISHED = 1
	k_CLOSED      = 2
)

type Connection interface {
	Send(m *Packet) error
	Close(err error)
	RemoteAddr() string
	LocalAddr() string
}

type TCPConnection struct {
	session     ISession
	conn        *net.TCPConn // 底层对象
	writeQueue  chan *Packet // 写队列
	readBuffer  *ByteBuffer  // 读缓冲
	writeBuffer *ByteBuffer  // 写缓冲
	readTimeout time.Duration
	state       int32
}

func (c *TCPConnection) init() {
	_ = c.conn.SetNoDelay(true)
	c.session.Bind(c) //
	c.setState(k_ESTABLISHED)
	c.session.HandleEvent(&SessionEvent{Session: c.session, Type: SESSION_NEW})
	go c.write()
	go c.read()
	log.Debug("[CONNECTION]: new connection", log.String("localAddr", c.LocalAddr()),
		log.String("remoteAddr", c.RemoteAddr()))
}

func (c *TCPConnection) Send(m *Packet) error {
	if c.getState() != k_ESTABLISHED {
		return errors.New("connection has already closed:" + c.String())
	}
	select {
	case c.writeQueue <- m:
		return nil
	default:
		return errors.New("connection write queue is full:" + c.String())
	}
}

func (c *TCPConnection) Close(err error) {
	if c.getState() == k_CLOSED {
		return
	}
	c.setState(k_CLOSED)
	close(c.writeQueue)
	c.session.HandleEvent(&SessionEvent{Session: c.session, Type: SESSION_CLOSED, Err: err})
	c.session = nil
}

//func (c *TCPConnection) Shutdown() {
//	if c.state == k_CLOSED {
//		return
//	}
//	c.state = k_CLOSED
//	c.session = nil
//	_ = c.conn.Close()
//	log.Debug("[CONNECTION]: connection shutdown", log.String("localAddr", c.LocalAddr()),
//		log.String("remoteAddr", c.RemoteAddr()))
//}

func (c *TCPConnection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *TCPConnection) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

func (c *TCPConnection) String() string {
	return c.LocalAddr() + "<--->" + c.RemoteAddr()
}

func (c *TCPConnection) handlePacket(p *Packet) {
	if c.getState() != k_ESTABLISHED { //可能已经被关闭
		return
	}
	c.session.HandleEvent(&SessionEvent{Session: c.session, Type: SESSION_PACKET, Pkt: p})
}

func (c *TCPConnection) handleError(err error) {
	if c.getState() != k_ESTABLISHED {
		return
	}
	c.Close(err)
}

func (c *TCPConnection) write() {
	for {
		// 阻塞组装包到write buffer
		if writeQueueClosed := c.packPacketBatch(); writeQueueClosed { // 此时写入通道被关闭了，关闭连接后直接退出
			_ = c.conn.Close()
			return
		}
		for c.writeBuffer.Length() > 0 { // 全部发出去
			err := c.writeBuffer.AppendWriter(c.conn) // 阻塞发送一次
			if err != nil {
				if e, ok := err.(net.Error); ok && e.Temporary() || e.Timeout() {
					continue
				}
				c.handleError(err)
				return
			}
		}
		c.writeBuffer.Reset()
	}
}

// 阻塞组装数据包到 write buffer中，返回write queue是不是阻塞，
// 保证尽量读取所有在队列的包，没有包时block住
func (c *TCPConnection) packPacketBatch() bool {
	var num = 0
	var queueSize = len(c.writeQueue) // 长度是0的时候，会阻塞，等待一个包触发，防止cpu空转
	var writeQueueClosed = true
	for pkt := range c.writeQueue {
		writeQueueClosed = false
		c.session.PreWriteHook(pkt) // 包预处理
		if pkt.Buff == nil {
			log.Debug("[CONNECTION]:write pkt buffer is nil,", log.Uint16("OpCode", pkt.GetOpCode()))
			continue
		}
		c.writeBuffer.Concat(pkt.Buff)
		num++
		// 保证有包的时候不会阻塞range 同时
		if num >= queueSize || c.writeBuffer.Space() < MAX_PACKET_LENGTH {
			break
		}
	}
	return writeQueueClosed
}

// 每次发送一个包
func (c *TCPConnection) packPacketOnce() bool {
	pkt := <-c.writeQueue
	if pkt == nil {
		return true
	}
	c.writeBuffer.Concat(pkt.Buff)
	return false
}

func (c *TCPConnection) read() {
	for {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
		err := c.readBuffer.AppendReader(c.conn) //阻塞读取一次
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				continue
			}
			c.handleError(err)
			return
		}
		for { //循环组包 直到pkt返回为nil，回到上面继续收包
			pkt, err := c.session.PreReadHook(c.readBuffer)
			if err != nil {
				c.handleError(err)
				return
			}
			if pkt == nil { //
				break
			}
			c.handlePacket(pkt)
		}
		if c.readBuffer.Space() < MAX_PACKET_LENGTH { //buffer剩余空间太小时，清理已读数据
			c.readBuffer.Slim()
		}
	}
}

func (c *TCPConnection) setState(state int) {
	atomic.StoreInt32(&c.state, int32(state))
}

func (c *TCPConnection) getState() int {
	return int(atomic.LoadInt32(&c.state))
}
