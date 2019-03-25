package net

import (
	"time"
)

type SessionEventType int

const (
	SESSION_NEW    SessionEventType = 0
	SESSION_CLOSED SessionEventType = 1
	SESSION_PACKET SessionEventType = 2
)

type SessionEvent struct {
	Type    SessionEventType
	Session ISession
	Err     error
	Pkt     *Packet
}

type ISession interface {
	PreReadHook(b *ByteBuffer) (*Packet, error) //io线程执行的解包函数
	PreWriteHook(p *Packet)                     //io线程执行的封包函数
	HandleEvent(event *SessionEvent)
	Bind(c Connection)
	Close(err error)
	RemoteAddr() string
	LocalAddr() string
}

func NewBaseSession(manager ISessionManager, sub SubSession) *BaseSession {
	return &BaseSession{manager: manager, SubSession: sub, heartBeatInterval: 5 * time.Second}
}

type SubSession interface {
	GetHeartBeatPacket() *Packet
	IsHeartBeatPacket(opCode uint16) bool
}

type BaseSession struct {
	SubSession
	conn              Connection
	manager           ISessionManager
	heartBeatTimer    *time.Timer
	heartBeatInterval time.Duration
	remoteAddr        string
	localAddr         string
}

func (s *BaseSession) startHeartBeat() {
	s.heartBeatTimer = time.AfterFunc(s.heartBeatInterval, func() {
		_ = s.Send(s.GetHeartBeatPacket())
		s.heartBeatTimer.Reset(s.heartBeatInterval)
	})
}

func (s *BaseSession) Bind(c Connection) {
	s.conn = c
	s.remoteAddr = s.conn.RemoteAddr()
	s.localAddr = s.conn.LocalAddr()
	s.startHeartBeat() //开启心跳包
}

func (s *BaseSession) Send(p *Packet) error {
	s.heartBeatTimer.Reset(s.heartBeatInterval) // 如果有新的包要发送，可以重置心跳包的发送间隔
	return s.conn.Send(p)
}

func (s *BaseSession) HandleEvent(event *SessionEvent) {
	// 过滤掉心跳包
	if event.Type == SESSION_PACKET && s.IsHeartBeatPacket(event.Pkt.GetOpCode()) {
		return
	}
	s.manager.HandleEvent(event)
}

//等待发送队列结束后关闭
func (s *BaseSession) Close(err error) {
	s.heartBeatTimer.Stop()
	s.conn.Close(err)
	s.conn = nil
}

func (s *BaseSession) RemoteAddr() string {
	return s.conn.RemoteAddr()
}

func (s *BaseSession) LocalAddr() string {
	return s.conn.LocalAddr()
}

type ISessionEventHandler interface {
	HandleNewSessionEvent(session ISession)
	HandleSessionClosedEvent(session ISession, err error)
	HandleSessionPacketEvent(session ISession, pkt *Packet)
}

type ISessionManager interface {
	Run()                                                 // 循环调用接口
	HandleEvent(event *SessionEvent)                      // 处理事件接口
	SetConnectFunc(cb func(addr string, ctx interface{})) // 设置重连回调,客户端连接才需要
	CreateSession(ctx interface{}) ISession               // 创建session
	Close(err error)                                      // 关闭
}

func NewBaseSessionManager(handler ISessionEventHandler) *BaseSessionManager {
	return &BaseSessionManager{
		ISessionEventHandler: handler,
		eventChan:            make(chan *SessionEvent, 1024*12),
	}
}

type BaseSessionManager struct {
	ISessionEventHandler
	eventChan chan *SessionEvent
}

func (m *BaseSessionManager) Run() {
	select {
	case event := <-m.eventChan:
		switch event.Type {
		case SESSION_NEW:
			m.HandleNewSessionEvent(event.Session)
		case SESSION_CLOSED:
			m.HandleSessionClosedEvent(event.Session, event.Err)
		case SESSION_PACKET:
			m.HandleSessionPacketEvent(event.Session, event.Pkt)
		}
	default:
	}
}

//设置重连函数 客户端连接才需要设置这个接口
func (m *BaseSessionManager) SetConnectFunc(cb func(addr string, ctx interface{})) {
}

func (m *BaseSessionManager) HandleEvent(event *SessionEvent) {
	m.eventChan <- event
}
