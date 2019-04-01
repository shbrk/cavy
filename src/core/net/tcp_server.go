package net

import (
	"core/log"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type ConnConfig struct {
	ReadBufferSize  int // 每个连接一个
	WriteBufferSize int //
	WriteQueueSize  int
	ReadTimeout     int // 读超时，单位s
}

type TCPServer struct {
	addr           string
	cfg            *ConnConfig
	ln             *net.TCPListener
	sessionManager ISessionManager
	IP             net.IP
	Port           int
	state          int32 // 原子变量
}

func NewTCPServer(addr string, cfg *ConnConfig, sessionManager ISessionManager) *TCPServer {
	return &TCPServer{cfg: cfg, addr: addr, sessionManager: sessionManager}
}

func (t *TCPServer) ListenAndServe() error {
	ln, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	t.Port = ln.Addr().(*net.TCPAddr).Port
	t.IP = ln.Addr().(*net.TCPAddr).IP
	t.ln = ln.(*net.TCPListener)
	go t.accept()
	return nil
}

func (t *TCPServer)Addr() string {
	return fmt.Sprintf("%s:%d",t.IP.String(),t.Port)
}

func (t *TCPServer) Close() {
	if t.IsClosed() {
		return
	}
	t.setClosed()
	_ = t.ln.Close()
	t.sessionManager.Close(nil)
}

func (t *TCPServer) IsClosed() bool {
	return atomic.LoadInt32(&t.state) == 1
}
func (t *TCPServer) setClosed() {
	atomic.StoreInt32(&t.state, 1)
}

func (t *TCPServer) GetSessionManager() ISessionManager {
	return t.sessionManager
}

func (t *TCPServer) accept() {
	var tempDelay time.Duration
	for {
		rw, e := t.ln.AcceptTCP()
		if e != nil {
			if t.IsClosed() {
				return
			}
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Error("[TCPServer]:Accept error", log.NamedError("err", e), log.Duration("retry in", tempDelay))
				time.Sleep(tempDelay)
				continue
			}
			log.Error("[TCPServer]:Accept error,stop working", log.NamedError("err", e))
			return
		}
		tempDelay = 0
		t.newConn(rw)
	}
}

func (t *TCPServer) newConn(conn *net.TCPConn) {
	NewTCPConnection(conn, t.sessionManager.CreateSession(), t.cfg.WriteQueueSize,
		t.cfg.ReadBufferSize, t.cfg.WriteBufferSize, t.cfg.ReadTimeout)
}
