package proto

import (
	"github.com/golang/protobuf/proto"

	"proto/client"
)

type ITest interface {
	GetID() int32
	GetCOUNT() int32
}

// 自动生成的消息解析
var OpCodeToMessage = map[uint16]func() proto.Message{
	uint16(client.OPCODE_C2S_PING): func() proto.Message { return &client.Ping{} },
	uint16(client.OPCODE_S2C_PING): func() proto.Message { return &client.Ping{} },
	uint16(client.OPCODE_TEST1):    func() proto.Message { return &client.Test1{} },
	uint16(client.OPCODE_TEST2):    func() proto.Message { return &client.Test2{} },
	uint16(client.OPCODE_TEST3):    func() proto.Message { return &client.Test3{} },
	//uint16(builtin.OPCODE_S2G_GS_REG): func() proto.Message { return &builtin.GSReg{} },
	//uint16(builtin.OPCODE_G2S_GS_REG): func() proto.Message { return &builtin.GSRegResult{} },
}
