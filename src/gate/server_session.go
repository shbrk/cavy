package gate

import (
	"core/log"
	"core/net"
	"github.com/golang/protobuf/proto"
	"sync/atomic"
)

// gate <-> server
//gate管理所有后端服务器的连接

func NewServerSession(id uint64, manager *ServerSessionManager) *ServerSession {
	session := &ServerSession{sessionID: id}
	baseSession := net.NewBaseSession(manager, session)
	session.BaseSession = baseSession
	return session
}

type ServerSession struct {
	*net.BaseSession
	sessionID uint64
}

func (s *ServerSession) PreReadHook(b *net.ByteBuffer) (*net.Packet, error) {
	return net.PreReadInternalHook(b)
}
func (s *ServerSession) PreWriteHook(p *net.Packet) {
	net.PreWriteInternalHook(p)
}

// 获取发送的心跳包
func (s *ServerSession) GetHeartBeatPacket() *net.Packet {
	return net.NewPacket(uint16(net.OP_G2S_HEART_BEAT), nil, 0)
}

//判断收到的包是不是心跳包
func (s *ServerSession) IsHeartBeatPacket(opCode uint16) bool {
	return opCode == uint16(net.OP_S2G_HEART_BEAT)
}

func (s *ServerSession) SendMsg(opCode uint16, guid uint64, msg proto.Message) {
	err := s.Send(net.NewPacket(opCode, msg, guid))
	if err != nil {
		log.Error("[SERVER_SESSION]：send message error", log.NamedError("err", err))
	}
}

func (s *ServerSession) SendData(opCode uint16, guid uint64, data []byte) {
	buf := net.NewByteBuffer(len(data))
	buf.Write(data)
	err := s.Send(net.NewPacketWithData(opCode, buf, guid))
	if err != nil {
		log.Error("[SERVER_SESSION]：send message error", log.NamedError("err", err))
	}
}

func NewServerSessionManager() *ServerSessionManager {
	manager := &ServerSessionManager{
		sessions: make(map[uint64]*ServerSession),
	}
	manager.BaseSessionManager = net.NewBaseSessionManager(manager)
	manager.ISessionEventHandler = manager
	return manager
}

type ServerSessionManager struct {
	*net.BaseSessionManager
	sessions        map[uint64]*ServerSession
	autoIncrementID uint64
}

func (m *ServerSessionManager) CreateSession(ctx interface{}) net.ISession {
	atomic.AddUint64(&m.autoIncrementID, 1)
	return NewServerSession(m.autoIncrementID, m)
}

func (m *ServerSessionManager) AddSession(session *ServerSession) {
	m.sessions[session.sessionID] = session
}
func (m *ServerSessionManager) RemoveSession(session *ServerSession) {
	delete(m.sessions, session.sessionID)
}

func (m *ServerSessionManager) HandleNewSessionEvent(session net.ISession) {
	serverSession, _ := session.(*ServerSession)
	log.Info("[SESSION]: new session ", log.String("localAddr", session.LocalAddr()),
		log.String("remoteAddr", session.RemoteAddr()))
	m.AddSession(serverSession)
}
func (m *ServerSessionManager) HandleSessionClosedEvent(session net.ISession, err error) {
	serverSession, _ := session.(*ServerSession)
	log.Info("[SESSION]: session closed", log.String("localAddr", session.LocalAddr()),
		log.String("remoteAddr", session.RemoteAddr()), log.NamedError("error", err))
	m.RemoveSession(serverSession)
}
func (m *ServerSessionManager) HandleSessionPacketEvent(session net.ISession, pkt *net.Packet) {
	serverSession, _ := session.(*ServerSession)
	log.Info(serverSession.RemoteAddr())
	//TODO 转发
}

func (m *ServerSessionManager) GetSession(id uint64) *ServerSession {
	return m.sessions[id]
}

func (m *ServerSessionManager) Close(err error) {
	for _, session := range m.sessions {
		session.Close(err)
	}
	m.sessions = make(map[uint64]*ServerSession)
}
