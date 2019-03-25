// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/client/opcode.proto

package client

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// 1-99 底层占用
// 100-19999是内部服务器消息
// 20000-49999客户端
// 50000-65535留给gm和其他节点使用
type OPCODE int32

const (
	//内置消息开始号
	OPCODE_EMPTY OPCODE = 0
	OPCODE_BEGIN OPCODE = 20000
	// Client发送到Server的Ping包
	OPCODE_C2S_PING OPCODE = 20001
	// Server发送回客户端的Ping
	OPCODE_S2C_PING OPCODE = 20002
	OPCODE_TEST1    OPCODE = 20003
	OPCODE_TEST2    OPCODE = 20004
	OPCODE_TEST3    OPCODE = 20005
	//内置消息结束号
	OPCODE_END OPCODE = 49999
)

var OPCODE_name = map[int32]string{
	0:     "EMPTY",
	20000: "BEGIN",
	20001: "C2S_PING",
	20002: "S2C_PING",
	20003: "TEST1",
	20004: "TEST2",
	20005: "TEST3",
	49999: "END",
}

var OPCODE_value = map[string]int32{
	"EMPTY":    0,
	"BEGIN":    20000,
	"C2S_PING": 20001,
	"S2C_PING": 20002,
	"TEST1":    20003,
	"TEST2":    20004,
	"TEST3":    20005,
	"END":      49999,
}

func (x OPCODE) String() string {
	return proto.EnumName(OPCODE_name, int32(x))
}

func (OPCODE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a99f3b373cd15c19, []int{0}
}

func init() {
	proto.RegisterEnum("client.OPCODE", OPCODE_name, OPCODE_value)
}

func init() { proto.RegisterFile("proto/client/opcode.proto", fileDescriptor_a99f3b373cd15c19) }

var fileDescriptor_a99f3b373cd15c19 = []byte{
	// 147 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2c, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xce, 0xc9, 0x4c, 0xcd, 0x2b, 0xd1, 0xcf, 0x2f, 0x48, 0xce, 0x4f, 0x49, 0xd5,
	0x03, 0x8b, 0x09, 0xb1, 0x41, 0x04, 0xb5, 0x8a, 0xb8, 0xd8, 0xfc, 0x03, 0x9c, 0xfd, 0x5d, 0x5c,
	0x85, 0x38, 0xb9, 0x58, 0x5d, 0x7d, 0x03, 0x42, 0x22, 0x05, 0x18, 0x84, 0xb8, 0xb9, 0x58, 0x9d,
	0x5c, 0xdd, 0x3d, 0xfd, 0x04, 0x16, 0xcc, 0x61, 0x14, 0xe2, 0xe3, 0xe2, 0x70, 0x36, 0x0a, 0x8e,
	0x0f, 0xf0, 0xf4, 0x73, 0x17, 0x58, 0x08, 0xe1, 0x07, 0x1b, 0x39, 0x43, 0xf8, 0x8b, 0xe6, 0x30,
	0x82, 0x14, 0x87, 0xb8, 0x06, 0x87, 0x18, 0x0a, 0x2c, 0x46, 0x70, 0x8c, 0x04, 0x96, 0x20, 0x38,
	0xc6, 0x02, 0x4b, 0xe7, 0x30, 0x0a, 0x71, 0x72, 0x31, 0xbb, 0xfa, 0xb9, 0x08, 0x9c, 0x6f, 0x63,
	0x4e, 0x62, 0x03, 0x3b, 0xc1, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xfb, 0xac, 0xa2, 0x68, 0x9f,
	0x00, 0x00, 0x00,
}
