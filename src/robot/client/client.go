package client

import (
	"core/log"
	"core/net"
	"github.com/golang/protobuf/proto"
	"proto/client"
	"sync"
	"time"
)

type Client struct {
	ID        int
	session   *ClientSession
	netClient *net.TCPClient
	count     uint64
	Wg        *sync.WaitGroup
}

func (c *Client) Connect(addr string) {
	c.netClient = net.NewTCPClient(5*time.Second, &net.ConnConfig{
		ReadBufferSize:  128 * 1024,
		WriteBufferSize: 128 * 1024,
		WriteQueueSize:  8 * 1024,
	}, NewClientSessionManager())
	session, err := c.netClient.SyncConnect(addr, time.Second*5, 1)
	if err != nil {
		log.Fatal("connect error", log.NamedError("err", err))
	}
	c.session = session.(*ClientSession)
}

func (c *Client) SendMsg(opCode uint16, msg proto.Message) {
	pkt := &net.Packet{IMsgHeader: &net.ClientHeaderIn{OpCode: opCode}, MsgBody: msg}
	err := c.session.Send(pkt)
	if err != nil {
		log.Error("send error", log.NamedError("err", err))
	}
}

func (c *Client) ping() {
	c.count += 1
	if c.count%3 == 1 {
		c.SendMsg(uint16(client.OPCODE_TEST1), &client.Test1{
			A: 1, B: "xxxaaa222", C: []uint64{1, 2, 3}, ID: int32(c.ID), COUNT: int32(c.count),
		})
	} else if c.count%3 == 2 {
		c.SendMsg(uint16(client.OPCODE_TEST2), &client.Test2{
			A: 222, Xxxaa: []int32{1, 2, 3}, Ccccxxxxxxxxx: 9, ID: int32(c.ID), COUNT: int32(c.count),
		})
	} else {
		c.SendMsg(uint16(client.OPCODE_TEST3), &client.Test3{
			Axxxxx: "axxx", Vxxxx: 2, A: 3, B: 1, C: 4, D: 2, F: "axx", ID: int32(c.ID), COUNT: int32(c.count),
		})
	}
	log.Info("", log.Int("id", c.ID), log.Uint64("count", c.count))
}

func (c *Client) Run() {
	for {
		c.netClient.Run()
		if c.count > 100 {
			c.Wg.Done()
			return
		}
		c.ping()
	}
}
