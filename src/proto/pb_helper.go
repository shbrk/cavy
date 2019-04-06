package proto

import (
	"github.com/golang/protobuf/proto"
	"proto/client"
)

// 自动生成的消息解析
var OpCodeToMessage = map[uint16]func() proto.Message{
	uint16(client.OPCODE_C2S_PING): func() proto.Message { return &client.Ping{} },
	uint16(client.OPCODE_S2C_PING): func() proto.Message { return &client.Ping{} },
}
