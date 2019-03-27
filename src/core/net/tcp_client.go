package net

import (
	"core/log"
	"net"
	"time"
)

type connectOp struct {
	addr    string
	session ISession
}

func NewTCPClient(timeout time.Duration, cfg *ConnConfig, sessionManager ISessionManager) *TCPClient {
	client := &TCPClient{timeout: timeout, cfg: cfg, sessionManager: sessionManager, connChan: make(chan *connectOp, 1024),
		closeChan: make(chan struct{}, 1)}
	client.init()
	return client
}

type TCPClient struct {
	cfg            *ConnConfig
	timeout        time.Duration
	sessionManager ISessionManager
	connChan       chan *connectOp
	closeChan      chan struct{}
}

func (t *TCPClient) init() {
	//为sessionManager绑定重连回调函数
	t.sessionManager.SetConnectFunc(func(addr string, session ISession) {
		t.Connect(addr, session)
	})
	go func() {
		for {
			co := <-t.connChan
			err := t.SyncConnect(co.addr, t.timeout, co.session)
			if err != nil {
				log.Error("[TCPClient]:tcp connect error", log.NamedError("err", err))
			}
		}
	}()
}

func (t *TCPClient) Run() {
	t.sessionManager.Run()
}
func (t *TCPClient) Close() {
	t.sessionManager.Close(nil)
}

func (t *TCPClient) GetSessionManager() ISessionManager {
	return t.sessionManager
}

func (t *TCPClient) Connect(addr string, session ISession) {
	t.connChan <- &connectOp{addr, session}
}

func (t *TCPClient) SyncConnect(addr string, timeout time.Duration, session ISession) error {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		return err
	}
	tcpConn, _ := conn.(*net.TCPConn)
	t.newConn(tcpConn, session)
	return nil
}

func (t *TCPClient) newConn(conn *net.TCPConn, session ISession) {
	NewTCPConnection(conn, session, t.cfg.WriteQueueSize, t.cfg.ReadBufferSize,
		t.cfg.WriteBufferSize, t.cfg.ReadTimeout)
}
