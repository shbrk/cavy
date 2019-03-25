package gate

import (
	"core/log"
	"core/net"
	"github.com/golang/protobuf/proto"
	"sync/atomic"
)

// gate 管理客户端的连接

func NewClientSession(id uint64, manager *ClientSessionManager) *ClientSession {
	session := &ClientSession{sessionID: id}
	baseSession := net.NewBaseSession(manager, session)
	session.BaseSession = baseSession
	return session
}

type ClientSession struct {
	*net.BaseSession
	sessionID uint64
	userID    uint64
}

func (c *ClientSession) PreReadHook(b *net.ByteBuffer) (*net.Packet, error) {
	return net.PreReadClientToGateHook(b)
}
func (c *ClientSession) PreWriteHook(p *net.Packet) {
	net.PreWriteGateToClientHook(p)
}

func (c *ClientSession) GetHeartBeatPacket() *net.Packet {
	return net.NewPacketToClient(uint16(net.OP_G2C_HEART_BEAT), nil)
}

func (c *ClientSession) IsHeartBeatPacket(opCode uint16) bool {
	return opCode == uint16(net.OP_C2G_HEART_BEAT)
}

func (c *ClientSession) SendMsg(opCode uint16, msg proto.Message) {
	err := c.Send(net.NewPacket(opCode,msg))
	if err != nil {
		log.Error("[CLIENT_SESSION]：send message error", log.NamedError("err", err))
	}
}

func (c *ClientSession) SendData(opCode uint16, guid uint64, data []byte) {
	buf := net.NewByteBuffer(len(data))
	buf.Write(data)
	err := c.Send(net.NewPackeToClienttWithData(opCode,buf))
	if err != nil {
		log.Error("[CLIENT_SESSION]：send message error", log.NamedError("err", err))
	}
}

func NewClientSessionManager() *ClientSessionManager {
	manager := &ClientSessionManager{
		sessions: make(map[uint64]*ClientSession),
	}
	manager.BaseSessionManager = net.NewBaseSessionManager(manager)
	manager.ISessionEventHandler = manager
	return manager
}

type ClientSessionManager struct {
	*net.BaseSessionManager
	sessions        map[uint64]*ClientSession
	autoIncrementID uint64
}

func (m *ClientSessionManager) CreateSession(ctx interface{}) net.ISession {
	atomic.AddUint64(&m.autoIncrementID, 1)
	return NewClientSession(m.autoIncrementID, m)
}

func (m *ClientSessionManager) AddSession(session *ClientSession) {
	m.sessions[session.sessionID] = session
}
func (m *ClientSessionManager) RemoveSession(session *ClientSession) {
	delete(m.sessions, session.sessionID)
}

func (m *ClientSessionManager) HandleNewSessionEvent(session net.ISession) {
	clientSession, _ := session.(*ClientSession)
	log.Info("[SESSION]: new session ", log.String("localAddr", session.LocalAddr()),
		log.String("remoteAddr", session.RemoteAddr()))
	m.AddSession(clientSession)
}
func (m *ClientSessionManager) HandleSessionClosedEvent(session net.ISession, err error) {
	clientSession, _ := session.(*ClientSession)
	log.Info("[SESSION]: session closed", log.String("localAddr", session.LocalAddr()),
		log.String("remoteAddr", session.RemoteAddr()), log.NamedError("error", err))
	clientSession.Close(err)
	m.RemoveSession(clientSession)
}
func (m *ClientSessionManager) HandleSessionPacketEvent(session net.ISession, pkt *net.Packet) {
	_, _ = session.(*ClientSession)
	//TODO 转发给后端服务器 测试代码
	//serverSession := gNodeGate.GetServerSession(1)
	//if serverSession == nil{
	//	log.Debug("can not found server session")
	//	return
	//}
	//m.test += 1
	//log.Debug("",log.Uint16("code",pkt.GetOpCode()),log.Uint64("guid",m.test))
	//serverSession.SendData(pkt.GetOpCode(), m.test, pkt.Buff.Data())
}

func (m *ClientSessionManager) GetSession(id uint64) *ClientSession {
	return m.sessions[id]
}

func (m *ClientSessionManager) Close(err error) {
	for _, session := range m.sessions {
		session.Close(err)
	}
	m.sessions = make(map[uint64]*ClientSession)
}
