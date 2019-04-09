package client

import (
	"core/log"
	"core/net"
	"fmt"
	"github.com/golang/protobuf/proto"
	"proto/client"
	"sync"
	"time"
)

type Client struct {
	ID             int
	session        *ClientSession
	netClient      *net.TCPClient
	sessionManager *ClientSessionManager
	count          uint64
	Wg             *sync.WaitGroup
}

func (c *Client) Connect(addr string) {
	c.sessionManager = NewClientSessionManager()
	c.sessionManager.SetPacketProcessor(c) // 设置包处理对象
	c.netClient = net.NewTCPClient(5*time.Second, &net.ConnConfig{
		ReadBufferSize:  128 * 1024,
		WriteBufferSize: 128 * 1024,
		WriteQueueSize:  8 * 1024,
	}, c.sessionManager)
	session := NewClientSession(uint64(c.ID), c.sessionManager)
	err := c.netClient.SyncConnect(addr, time.Second*5, session)
	if err != nil {
		log.Fatal("connect error", log.NamedError("err", err))
	}
	c.session = session

}

func (c *Client) SendMsg(opCode uint16, msg proto.Message) {
	pkt := &net.Packet{IMsgHeader: &net.ClientHeaderIn{OpCode: opCode}, MsgBody: msg}
	err := c.session.Send(pkt)
	if err != nil {
		log.Error("send error", log.NamedError("err", err))
	}
}

func (c *Client) ProcessPacket(session net.ISession, pkt *net.Packet) {
	if pkt.MsgBody == nil {
		log.Warn("[CLIENT]:msg body not auto unmarshal", log.Uint16("opCode", pkt.GetOpCode()))
		return
	}
	var opCode = client.OPCODE(pkt.GetOpCode())
	switch opCode {
	case client.OPCODE_S2C_LOGIN:
		c.onLoginBack(pkt.MsgBody.(*client.S2C_Login))
	}
}

func (c *Client) login() {
	var token = fmt.Sprintf("robot%d", c.ID)
	var msg = &client.C2S_Login{Token: token, AreaID: 1}
	c.SendMsg(uint16(client.OPCODE_C2S_LOGIN), msg)
}

func (c *Client) onLoginBack(msg *client.S2C_Login) {
	log.Info("login success")
}

func (c *Client) Run() {
	for {
		event := <-c.sessionManager.EventChan
		c.sessionManager.HandleEvent(event)
	}
}
