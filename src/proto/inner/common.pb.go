// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/inner/common.proto

package inner

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

// GS向Gate注册
// S2G_GS_REG = 7;
type GSReg struct {
	AreaId               int32    `protobuf:"varint,1,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GSReg) Reset()         { *m = GSReg{} }
func (m *GSReg) String() string { return proto.CompactTextString(m) }
func (*GSReg) ProtoMessage()    {}
func (*GSReg) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd220dc4ffc9a07c, []int{0}
}

func (m *GSReg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GSReg.Unmarshal(m, b)
}
func (m *GSReg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GSReg.Marshal(b, m, deterministic)
}
func (m *GSReg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GSReg.Merge(m, src)
}
func (m *GSReg) XXX_Size() int {
	return xxx_messageInfo_GSReg.Size(m)
}
func (m *GSReg) XXX_DiscardUnknown() {
	xxx_messageInfo_GSReg.DiscardUnknown(m)
}

var xxx_messageInfo_GSReg proto.InternalMessageInfo

func (m *GSReg) GetAreaId() int32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

// Gate返回给GS的注册结果
// G2S_GS_REG = 8;
type GSRegResult struct {
	Ret                  bool     `protobuf:"varint,1,opt,name=ret,proto3" json:"ret,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GSRegResult) Reset()         { *m = GSRegResult{} }
func (m *GSRegResult) String() string { return proto.CompactTextString(m) }
func (*GSRegResult) ProtoMessage()    {}
func (*GSRegResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd220dc4ffc9a07c, []int{1}
}

func (m *GSRegResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GSRegResult.Unmarshal(m, b)
}
func (m *GSRegResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GSRegResult.Marshal(b, m, deterministic)
}
func (m *GSRegResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GSRegResult.Merge(m, src)
}
func (m *GSRegResult) XXX_Size() int {
	return xxx_messageInfo_GSRegResult.Size(m)
}
func (m *GSRegResult) XXX_DiscardUnknown() {
	xxx_messageInfo_GSRegResult.DiscardUnknown(m)
}

var xxx_messageInfo_GSRegResult proto.InternalMessageInfo

func (m *GSRegResult) GetRet() bool {
	if m != nil {
		return m.Ret
	}
	return false
}

func init() {
	proto.RegisterType((*GSReg)(nil), "inner.GSReg")
	proto.RegisterType((*GSRegResult)(nil), "inner.GSRegResult")
}

func init() { proto.RegisterFile("proto/inner/common.proto", fileDescriptor_cd220dc4ffc9a07c) }

var fileDescriptor_cd220dc4ffc9a07c = []byte{
	// 112 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x28, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0xcf, 0xcc, 0xcb, 0x4b, 0x2d, 0xd2, 0x4f, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0xd3, 0x03,
	0x0b, 0x09, 0xb1, 0x82, 0xc5, 0x94, 0x14, 0xb8, 0x58, 0xdd, 0x83, 0x83, 0x52, 0xd3, 0x85, 0xc4,
	0xb9, 0xd8, 0x13, 0x8b, 0x52, 0x13, 0xe3, 0x33, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83,
	0xd8, 0x40, 0x5c, 0xcf, 0x14, 0x25, 0x79, 0x2e, 0x6e, 0xb0, 0x8a, 0xa0, 0xd4, 0xe2, 0xd2, 0x9c,
	0x12, 0x21, 0x01, 0x2e, 0xe6, 0xa2, 0xd4, 0x12, 0xb0, 0x1a, 0x8e, 0x20, 0x10, 0x33, 0x89, 0x0d,
	0x6c, 0xa0, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xa3, 0x9c, 0x9e, 0x90, 0x6c, 0x00, 0x00, 0x00,
}
