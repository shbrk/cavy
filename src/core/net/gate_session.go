package net

import (
	"core/log"
	"errors"
	"github.com/golang/protobuf/proto"
	"proto/inner"
	"sync/atomic"
	"time"
)

// server ----> gate
//管理所有后端服务器对gate的连接

const MAX_RECONNECT_COUNT = 20

func NewGateSession(bootID int, sessionID uint64, manager *GateSessionManager) *GateSession {
	session := &GateSession{sessionID: sessionID, bootID: bootID}
	baseSession := NewBaseSession(manager, session)
	session.BaseSession = baseSession
	return session
}

type GateSession struct {
	*BaseSession
	sessionID      uint64
	bootID         int
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
		log.Error("[GATE_SESSION]：send message error", log.NamedError("err", err))
	}
}

func (g *GateSession) RegToGate(opCode inner.OPCODE, msg proto.Message) {
	g.SendMsg(uint16(opCode), msg)
}
func (g *GateSession) OnRegBack(ret bool) {
	if ret {
		log.Info("[GATE_SESSION] reg to gate ok", log.Int("bootID", g.bootID),
			log.String("localAddr", g.LocalAddr()),
			log.String("remoteAddr", g.RemoteAddr()))
	} else {
		log.Error("[GATE_SESSION] reg to gate failed", log.Int("bootID", g.bootID),
			log.String("localAddr", g.LocalAddr()),
			log.String("remoteAddr", g.RemoteAddr()))
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
	*BaseSessionManager                                     //基类
	sessions            map[uint64]*GateSession             //以sessionID为key的map
	sessionsWithBootID  map[int]*GateSession                //以bootID为key的map
	connectFunc         func(addr string, session ISession) //连接回调
	reconnectInterval   time.Duration                       // 重连最小周期
	autoIncrementID     uint64                              //自增ID
	regOpCode           inner.OPCODE                     //向gate注册的消息ID
	regMsg              proto.Message                       // 向gate注册的消息
}

func (m *GateSessionManager) CreateSession() ISession {
	return nil
}

func (m *GateSessionManager) CreateSessionWithBootID(bootID int) ISession {
	atomic.AddUint64(&m.autoIncrementID, 1)
	return NewGateSession(bootID, m.autoIncrementID, m)
}

func (m *GateSessionManager) AddSession(session *GateSession) {
	m.sessions[session.sessionID] = session
	m.sessionsWithBootID[session.bootID] = session
}
func (m *GateSessionManager) RemoveSession(session *GateSession) {
	delete(m.sessions, session.sessionID)
	delete(m.sessionsWithBootID, session.bootID)
}

func (m *GateSessionManager) SetRegisterInfo(opCode inner.OPCODE, msg proto.Message) {
	m.regOpCode = opCode
	m.regMsg = msg
}

func (m *GateSessionManager) SetConnectFunc(cb func(addr string, session ISession)) {
	m.connectFunc = cb
}

func (m *GateSessionManager) HandleEtcdEventPut(bootID int, addr string) {
	m.connectFunc(addr, m.CreateSessionWithBootID(bootID))
}

func (m *GateSessionManager) HandleEtcdEventDelete(booID int) {
	session := m.GetSessionByBootID(booID)
	if session != nil {
		session.Close(errors.New("etcd node is down"))
		m.RemoveSession(session)
	}
	log.Error("[GATE_SESSION]: etcd node is down", log.Int("bootID", booID))
}

func (m *GateSessionManager) HandleNewSessionEvent(session ISession) {
	gateSession, _ := session.(*GateSession)
	if oldSession := m.GetSession(gateSession.sessionID); oldSession != nil {
		if oldSession.isReconnecting {
			oldSession.StopReconnect()
		}
		log.Info("[GATE_SESSION]：reconnection successful", log.String("localAddr", session.LocalAddr()),
			log.String("remoteAddr", session.RemoteAddr()))

	} else {
		if oldSession := m.GetSessionByBootID(gateSession.bootID); oldSession != nil {
			oldSession.Close(errors.New("[GATE_SESSION]: new session has duplicate id"))
			m.RemoveSession(oldSession)
		}
		log.Info("[GATE_SESSION]: new session ", log.String("localAddr", session.LocalAddr()),
			log.String("remoteAddr", session.RemoteAddr()))
		m.AddSession(gateSession)
		gateSession.RegToGate(m.regOpCode, m.regMsg)
	}
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
			m.RemoveSession(gateSession)
			log.Error("[GATE_SESSION]: auto reconnect failed", log.Int("tryCount", count))
		})
	}
}

func (m *GateSessionManager) HandleSessionPacketEvent(session ISession, pkt *Packet) {
	_, _ = session.(*GateSession)
	if m.Processor != nil {
		m.Processor.ProcessPacket(session, pkt)
	}
}

func (m *GateSessionManager) GetSession(id uint64) *GateSession {
	return m.sessions[id]
}

func (m *GateSessionManager) GetSessionByBootID(bootID int) *GateSession {
	return m.sessionsWithBootID[bootID]
}

func (m *GateSessionManager) Close(err error) {
	for _, session := range m.sessions {
		session.Close(err)
	}
	m.sessions = make(map[uint64]*GateSession)
}
