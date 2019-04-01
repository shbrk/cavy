package net

import (
	"core/log"
	"errors"
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
	PostEvent(event *SessionEvent)
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
	if s.conn != nil {
		s.conn.Close(errors.New("rebind new connection,so close old connection "))
	}
	s.conn = c
	s.remoteAddr = s.conn.RemoteAddr()
	s.localAddr = s.conn.LocalAddr()
	s.startHeartBeat() //开启心跳包
}

func (s *BaseSession) Send(p *Packet) error {
	s.heartBeatTimer.Reset(s.heartBeatInterval) // 如果有新的包要发送，可以重置心跳包的发送间隔
	if s.conn != nil {
		return s.conn.Send(p)
	}
	return nil
}

func (s *BaseSession) PostEvent(event *SessionEvent) {
	// 过滤掉心跳包
	if event.Type == SESSION_PACKET && s.IsHeartBeatPacket(event.Pkt.GetOpCode()) {
		return
	}
	s.manager.PostEvent(event)
}

//等待发送队列结束后关闭
func (s *BaseSession) Close(err error) {
	s.heartBeatTimer.Stop()
	s.conn.Close(err)
	s.conn = nil
}

func (s *BaseSession) RemoteAddr() string {
	return s.remoteAddr
}

func (s *BaseSession) LocalAddr() string {
	return s.localAddr
}

type ISessionEventHandler interface {
	HandleNewSessionEvent(session ISession)
	HandleSessionClosedEvent(session ISession, err error)
	HandleSessionPacketEvent(session ISession, pkt *Packet)
}

type ISessionManager interface {
	HandleEvent(event *SessionEvent)                       // 处理事件接口
	PostEvent(event *SessionEvent)                         // 投递事件接口
	SetConnectFunc(cb func(addr string, session ISession)) // 设置重连回调,客户端连接才需要
	CreateSession() ISession                               // 创建session
	Close(err error)                                       // 关闭
}

func NewBaseSessionManager(handler ISessionEventHandler) *BaseSessionManager {
	return &BaseSessionManager{
		ISessionEventHandler: handler,
		EventChan:            make(chan *SessionEvent, 1024*12),
	}
}

type BaseSessionManager struct {
	ISessionEventHandler
	EventChan chan *SessionEvent
}

func (m *BaseSessionManager) HandleEvent(event *SessionEvent) {
	if event == nil {
		log.Error("[SESSION]: event is nil")
		return
	}
	switch event.Type {
	case SESSION_NEW:
		m.HandleNewSessionEvent(event.Session)
	case SESSION_CLOSED:
		m.HandleSessionClosedEvent(event.Session, event.Err)
	case SESSION_PACKET:
		m.HandleSessionPacketEvent(event.Session, event.Pkt)
	}
}

//设置重连函数 客户端连接才需要设置这个接口
func (m *BaseSessionManager) SetConnectFunc(cb func(addr string, session ISession)) {
}

func (m *BaseSessionManager) PostEvent(event *SessionEvent) {
	m.EventChan <- event
}
