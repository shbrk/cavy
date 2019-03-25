package net

import (
	"core/log"
	"errors"
	"github.com/golang/protobuf/proto"
	"time"
)

// server ----> gate
//管理所有后端服务器对gate的连接

func NewGateSession(id uint64, manager *GateSessionManager) *GateSession {
	session := &GateSession{sessionID: id}
	baseSession := NewBaseSession(manager, session)
	session.BaseSession = baseSession
	return session
}

type GateSession struct {
	*BaseSession
	sessionID      uint64
	reconnectTimer *time.Timer
	isReconnecting bool // 是否正在重连中
}

func (g *GateSession) PreReadHook(b *ByteBuffer) (*Packet, error) {
	return PreReadInternalHook(b)
}

func (g *GateSession) PreWriteHook(p *Packet) {
	PreWriteInternalHook(p)
}

func (g *GateSession) StartReconnect(interval time.Duration, cb func(addr string)) {
	g.isReconnecting = true
	var nextInterval = interval
	var reconnectCount = 0
	g.reconnectTimer = time.AfterFunc(nextInterval, func() {
		nextInterval = nextInterval * 2
		reconnectCount += 1
		g.reconnectTimer.Reset(nextInterval)
		cb(g.remoteAddr)
	})
}

func (g *GateSession) Close(err error) {
	if g.reconnectTimer != nil {
		g.reconnectTimer.Stop()
	}
	g.BaseSession.Close(err)
}

func (g *GateSession) GetHeartBeatPacket() *Packet {
	return NewPacket(uint16(OP_S2G_HEART_BEAT), nil,0)
}
func (g *GateSession) IsHeartBeatPacket(opCode uint16) bool {
	return opCode == uint16(OP_G2S_HEART_BEAT)
}
func (g *GateSession) SendMsg(opCode uint16, msg proto.Message) {
	err := g.Send(&Packet{IMsgHeader: &InternalHeader{OpCode: opCode}, MsgBody: msg})
	if err != nil {
		log.Error("[GATE_SESSION]：send message error", log.NamedError("err", err))
	}
}

func (g *GateSession) SendData(opCode uint16, guid uint64, data []byte) {
	buf := NewByteBuffer(len(data))
	buf.Write(data)
	err := g.Send(&Packet{IMsgHeader: &InternalHeader{OpCode: opCode}, Buff: buf})
	if err != nil {
		log.Error("[CLIENT_SESSION]：send message error", log.NamedError("err", err))
	}
}

func NewGateSessionManager() *GateSessionManager {
	manager := &GateSessionManager{
		sessions:          make(map[uint64]*GateSession),
		reconnectInterval: 5 * time.Second,
	}
	manager.BaseSessionManager = NewBaseSessionManager(manager)
	manager.ISessionEventHandler = manager
	return manager
}

type ISessionConsumer interface {
	OnSessionMessage(sessionID, guid string, opCode uint16, msg proto.Message)
	OnSessionCreate(sessionID uint64)
	OnSessionClose(sessionID uint64)
	SetSendMessageFunc(func(sessionID uint64, guid string, opCode uint16, msg proto.Message))
	SetSendDataFunc(func(sessionID uint64, guid string, opCode uint16, data []byte))
}

type GateSessionManager struct {
	*BaseSessionManager
	sessions          map[uint64]*GateSession
	connectFunc       func(addr string, ctx interface{})
	reconnectInterval time.Duration
}

func (m *GateSessionManager) CreateSession(ctx interface{}) ISession {
	id := ctx.(int)
	return NewGateSession(uint64(id), m)
}

func (m *GateSessionManager) AddSession(session *GateSession) {
	m.sessions[session.sessionID] = session
}
func (m *GateSessionManager) RemoveSession(session *GateSession) {
	delete(m.sessions, session.sessionID)
}

func (m *GateSessionManager) SetConnectFunc(cb func(addr string, ctx interface{})) {
	m.connectFunc = cb
}

func (m *GateSessionManager) HandleEtcdEventPut(bootID int, addr string) {
	if session := m.GetSession(uint64(bootID)); session != nil {
		log.Error("[SESSION]:etcd new put has duplicate id", log.Int("bootID", bootID))
		return
	}
	m.connectFunc(addr, bootID)
}

func (m *GateSessionManager) HandleEtcdEventDelete(booID int) {
	log.Warn("[SESSION]: etcd node is down", log.Int("bootID", booID))
}

func (m *GateSessionManager) HandleNewSessionEvent(session ISession) {
	gateSession, _ := session.(*GateSession)
	if oldSession := m.GetSession(gateSession.sessionID); oldSession != nil {
		if oldSession.isReconnecting {
			oldSession.Close(errors.New("reconnect new session success"))
		} else {
			gateSession.Close(errors.New("[SESSION]: new session has duplicate id"))
			log.Error("[SESSION]: new session has duplicate id", log.Uint64("sessionID", gateSession.sessionID))
			return
		}
	} else {
		log.Info("[SESSION]: new session ", log.String("localAddr", session.LocalAddr()),
			log.String("remoteAddr", session.RemoteAddr()))
	}
	m.AddSession(gateSession)
}
func (m *GateSessionManager) HandleSessionClosedEvent(session ISession, err error) {
	gateSession, _ := session.(*GateSession)
	log.Info("[SESSION]: session closed", log.String("localAddr", session.LocalAddr()),
		log.String("remoteAddr", session.RemoteAddr()), log.NamedError("error", err))
	if m.GetSession(gateSession.sessionID) != nil { //如果逻辑层没有移除session 自动重连
		gateSession.StartReconnect(m.reconnectInterval, func(addr string) {
			m.connectFunc(addr, int(gateSession.sessionID))
		})
	}
}

func (m *GateSessionManager) HandleSessionPacketEvent(session ISession, pkt *Packet) {
	_, _ = session.(*GateSession)
	//msg := pkt.MsgBody.(proto2.ITest)
	//head := pkt.IMsgHeader.(*InternalHeader)
	// log.Debug("", log.Uint16("code", head.GetOpCode()), log.Uint64("guid", head.Guid), log.Int32("id", msg.GetID()), log.Int32("count", msg.GetCOUNT()))
	//TODO 派发给逻辑层
	//	log.Debug("msg",log.Uint64("GateBackRev",msg.GateBackRev),log.Uint64("GsRecv",msg.GsRecv))
}

func (m *GateSessionManager) GetSession(id uint64) *GateSession {
	return m.sessions[id]
}

func (m *GateSessionManager) Close(err error) {
	for _, session := range m.sessions {
		session.Close(err)
	}
	m.sessions = make(map[uint64]*GateSession)
}
