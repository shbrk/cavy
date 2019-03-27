package client

import (
	"core/log"
	"core/net"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"sync/atomic"
)

func NewClientSession(id uint64, manager *ClientSessionManager) *ClientSession {
	session := &ClientSession{sessionID: id}
	baseSession := net.NewBaseSession(manager, session)
	session.BaseSession = baseSession
	return session
}

type ClientSession struct {
	*net.BaseSession
	sessionID uint64
}

func (c *ClientSession) PreReadHook(b *net.ByteBuffer) (*net.Packet, error) {
	head := &net.ClientHeaderOut{}
	buffSize := b.Length()
	if buffSize < head.GetLength() {
		return nil, nil
	}
	head.PeekSize(b)
	if head.Size > net.MAX_PACKET_LENGTH {
		return nil, errors.New(fmt.Sprintf("packet length is to large:%d", head.Size))
	}
	if uint32(buffSize) < head.Size {
		return nil, nil
	}
	head.Decode(b)
	pkt := &net.Packet{IMsgHeader: head}
	bodySize := int(head.Size) - head.GetLength()
	if bodySize > 0 {
		pkt.Buff = b.Copy(bodySize) // 直接从copy到buffer中返回
	}
	return pkt, nil
}
func (c *ClientSession) PreWriteHook(p *net.Packet) {
	head := p.IMsgHeader.(*net.ClientHeaderIn)
	if p.MsgBody != nil { // Gate自己构建的包
		buff, err := proto.Marshal(p.MsgBody.(proto.Message))
		if err != nil {
			log.Error("[PACKET]:marshal message error", log.Uint16("OpCode", head.GetOpCode()))
		}
		head.Size = uint32(head.GetLength() + len(buff))
		p.Buff = net.NewByteBuffer(int(head.Size))
		head.Encode(p.Buff)
		p.Buff.Write(buff)
	} else { // Gate转发的包 或者是空包
		if p.Buff != nil {
			size := head.GetLength() + p.Buff.Length()
			head.Size = uint32(size)
			buf := net.NewByteBuffer(int(head.Size))
			head.Encode(buf)
			buf.Concat(p.Buff)
			p.Buff = buf
		} else { // 空包
			head.Size = uint32(head.GetLength())
			p.Buff = net.NewByteBuffer(int(head.Size))
			head.Encode(p.Buff)
		}
	}
}

func (c *ClientSession) GetHeartBeatPacket() *net.Packet {
	return &net.Packet{IMsgHeader: &net.ClientHeaderIn{OpCode: net.OP_C2G_HEART_BEAT}}
}

func (c *ClientSession) IsHeartBeatPacket(opCode uint16) bool {
	return opCode == uint16(net.OP_G2C_HEART_BEAT)
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

func (m *ClientSessionManager) CreateSession() net.ISession {
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
	m.RemoveSession(clientSession)
}
func (m *ClientSessionManager) HandleSessionPacketEvent(session net.ISession, pkt *net.Packet) {
	clientSession, _ := session.(*ClientSession)
	log.Info(clientSession.RemoteAddr())
	//TODO 转发给后端服务器
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
