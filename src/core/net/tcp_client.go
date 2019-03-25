package net

import (
	"core/log"
	"net"
	"time"
)

type connectOp struct {
	addr string
	ctx  interface{}
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
	t.sessionManager.SetConnectFunc(func(addr string, ctx interface{}) {
		t.Connect(addr, ctx)
	})
	go func() {
		for {
			co := <-t.connChan
			_, err := t.SyncConnect(co.addr, t.timeout, co.ctx)
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

func (t *TCPClient) Connect(addr string, ctx interface{}) {
	t.connChan <- &connectOp{addr, ctx}
}

func (t *TCPClient) SyncConnect(addr string, timeout time.Duration, ctx interface{}) (ISession, error) {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcpConn, _ := conn.(*net.TCPConn)
	return t.newConn(tcpConn, ctx), nil
}

func (t *TCPClient) newConn(conn *net.TCPConn, ctx interface{}) ISession {
	session := t.sessionManager.CreateSession(ctx)
	NewTCPConnection(conn, session, t.cfg.WriteQueueSize,
		t.cfg.ReadBufferSize, t.cfg.WriteBufferSize, t.cfg.ReadTimeout)
	return session
}
