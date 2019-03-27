package net

import (
	"core/log"
	"errors"
	"github.com/golang/protobuf/proto"
	"time"
)

// server ----> gate
//管理所有后端服务器对gate的连接

const MAX_RECONNECT_COUNT = 20

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

func (g *GateSession) StartReconnect(interval time.Duration, cb func(addr string), fail func(count int)) {
	g.isReconnecting = true
	var nextInterval = interval
	var reconnectCount = 0
	g.reconnectTimer = time.AfterFunc(nextInterval, func() {
		nextInterval = nextInterval * 2
		reconnectCount += 1
		if reconnectCount > MAX_RECONNECT_COUNT {
			g.StopReconnect()
			fail(reconnectCount)
		}
		g.reconnectTimer.Reset(nextInterval)
		cb(g.remoteAddr)
	})
}

func (g *GateSession) StopReconnect() {
	if !g.isReconnecting {
		return
	}
	if g.reconnectTimer != nil {
		g.reconnectTimer.Stop()
		g.reconnectTimer = nil
	}
}

func (g *GateSession) Close(err error) {
	if g.reconnectTimer != nil {
		g.reconnectTimer.Stop()
	}
	g.BaseSession.Close(err)
}

func (g *GateSession) GetHeartBeatPacket() *Packet {
	return NewPacket(uint16(OP_S2G_HEART_BEAT), nil, 0)
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
	connectFunc       func(addr string, session ISession)
	reconnectInterval time.Duration
}

func (m *GateSessionManager) CreateSession() ISession {
	return nil
}

func (m *GateSessionManager) AddSession(session *GateSession) {
	m.sessions[session.sessionID] = session
}
func (m *GateSessionManager) RemoveSession(session *GateSession) {
	delete(m.sessions, session.sessionID)
}

func (m *GateSessionManager) SetConnectFunc(cb func(addr string, session ISession)) {
	m.connectFunc = cb
}

func (m *GateSessionManager) HandleEtcdEventPut(bootID int, addr string) {
	//if session := m.GetSession(uint64(bootID)); session != nil {
	//	log.Error("[GATE_SESSION]:etcd new put has duplicate id", log.Int("bootID", bootID))
	//	return
	//}
	m.connectFunc(addr, NewGateSession(uint64(bootID), m))
}

func (m *GateSessionManager) HandleEtcdEventDelete(booID int) {
	log.Warn("[GATE_SESSION]: etcd node is down", log.Int("bootID", booID))
	//if session := m.GetSession(uint64(booID)); session != nil {
	//	session.Close(errors.New("etcd node is down"))
	//	m.RemoveSession(session)
	//}
}

func (m *GateSessionManager) HandleNewSessionEvent(session ISession) {
	gateSession, _ := session.(*GateSession)
	if oldSession := m.GetSession(gateSession.sessionID); oldSession != nil {
		if oldSession == gateSession { // 已经存在的session
			if oldSession.isReconnecting {
				log.Info("[GATE_SESSION]：reconnect succeed", log.String("localAddr", session.LocalAddr()),
					log.String("remoteAddr", session.RemoteAddr()))
				oldSession.StopReconnect()
			} else { // 可能是延迟比较大 多次重连结果都连接成功
				log.Warn("reconnect maybe to make connection twice", log.String("localAddr", session.LocalAddr()),
					log.String("remoteAddr", session.RemoteAddr()))
			}
		} else {
			if oldSession.isReconnecting { // 可能是对方重启了ID没变 addr变了
				oldSession.Close(errors.New("[GATE_SESSION]:new session has duplicate id,may same remote changed addr"))
			} else {
				gateSession.Close(errors.New("[GATE_SESSION]: new session has duplicate id"))
				log.Error("[GATE_SESSION]: new session has duplicate id", log.Uint64("sessionID", gateSession.sessionID))
				return
			}
		}
	} else {
		log.Info("[GATE_SESSION]: new session ", log.String("localAddr", session.LocalAddr()),
			log.String("remoteAddr", session.RemoteAddr()))
	}
	m.AddSession(gateSession)
}
func (m *GateSessionManager) HandleSessionClosedEvent(session ISession, err error) {
	gateSession, _ := session.(*GateSession)
	log.Info("[GATE_SESSION]: session closed", log.String("localAddr", session.LocalAddr()),
		log.String("remoteAddr", session.RemoteAddr()), log.NamedError("error", err))
	if m.GetSession(gateSession.sessionID) != nil { //如果逻辑层没有移除session 自动重连
		gateSession.StartReconnect(m.reconnectInterval, func(addr string) {
			m.connectFunc(addr, gateSession)
			log.Info("[GATE_SESSION]:try reconnect :" + addr)
		}, func(count int) {
			log.Error("[GATE_SESSION]: auto reconnect failed", log.Int("tryCount", count))
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
