// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/client/common.proto

package client

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	entity "proto/entity"
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

// ping包
// C2S_PING = 5;S2C_PING = 6;
type Ping struct {
	ClientSend           uint64   `protobuf:"fixed64,1,opt,name=client_send,json=clientSend,proto3" json:"client_send,omitempty"`
	GateRecv             uint64   `protobuf:"fixed64,2,opt,name=gate_recv,json=gateRecv,proto3" json:"gate_recv,omitempty"`
	GateSend             uint64   `protobuf:"fixed64,3,opt,name=gate_send,json=gateSend,proto3" json:"gate_send,omitempty"`
	GsRecv               uint64   `protobuf:"fixed64,4,opt,name=gs_recv,json=gsRecv,proto3" json:"gs_recv,omitempty"`
	GsSend               uint64   `protobuf:"fixed64,5,opt,name=gs_send,json=gsSend,proto3" json:"gs_send,omitempty"`
	GateBackRev          uint64   `protobuf:"fixed64,6,opt,name=gate_back_rev,json=gateBackRev,proto3" json:"gate_back_rev,omitempty"`
	GateBackSend         uint64   `protobuf:"fixed64,7,opt,name=gate_back_send,json=gateBackSend,proto3" json:"gate_back_send,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ping) Reset()         { *m = Ping{} }
func (m *Ping) String() string { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()    {}
func (*Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae54b1244e5266d5, []int{0}
}

func (m *Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ping.Unmarshal(m, b)
}
func (m *Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ping.Marshal(b, m, deterministic)
}
func (m *Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ping.Merge(m, src)
}
func (m *Ping) XXX_Size() int {
	return xxx_messageInfo_Ping.Size(m)
}
func (m *Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Ping proto.InternalMessageInfo

func (m *Ping) GetClientSend() uint64 {
	if m != nil {
		return m.ClientSend
	}
	return 0
}

func (m *Ping) GetGateRecv() uint64 {
	if m != nil {
		return m.GateRecv
	}
	return 0
}

func (m *Ping) GetGateSend() uint64 {
	if m != nil {
		return m.GateSend
	}
	return 0
}

func (m *Ping) GetGsRecv() uint64 {
	if m != nil {
		return m.GsRecv
	}
	return 0
}

func (m *Ping) GetGsSend() uint64 {
	if m != nil {
		return m.GsSend
	}
	return 0
}

func (m *Ping) GetGateBackRev() uint64 {
	if m != nil {
		return m.GateBackRev
	}
	return 0
}

func (m *Ping) GetGateBackSend() uint64 {
	if m != nil {
		return m.GateBackSend
	}
	return 0
}

type EntityContainer struct {
	EntityType           entity.TYPE `protobuf:"varint,1,opt,name=entityType,proto3,enum=entity.TYPE" json:"entityType,omitempty"`
	List                 [][]byte    `protobuf:"bytes,2,rep,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *EntityContainer) Reset()         { *m = EntityContainer{} }
func (m *EntityContainer) String() string { return proto.CompactTextString(m) }
func (*EntityContainer) ProtoMessage()    {}
func (*EntityContainer) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae54b1244e5266d5, []int{1}
}

func (m *EntityContainer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntityContainer.Unmarshal(m, b)
}
func (m *EntityContainer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntityContainer.Marshal(b, m, deterministic)
}
func (m *EntityContainer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntityContainer.Merge(m, src)
}
func (m *EntityContainer) XXX_Size() int {
	return xxx_messageInfo_EntityContainer.Size(m)
}
func (m *EntityContainer) XXX_DiscardUnknown() {
	xxx_messageInfo_EntityContainer.DiscardUnknown(m)
}

var xxx_messageInfo_EntityContainer proto.InternalMessageInfo

func (m *EntityContainer) GetEntityType() entity.TYPE {
	if m != nil {
		return m.EntityType
	}
	return entity.TYPE_None
}

func (m *EntityContainer) GetList() [][]byte {
	if m != nil {
		return m.List
	}
	return nil
}

// 发送某些实体的所有列表
type EntityList struct {
	Owner                uint64             `protobuf:"varint,1,opt,name=owner,proto3" json:"owner,omitempty"`
	List                 []*EntityContainer `protobuf:"bytes,2,rep,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *EntityList) Reset()         { *m = EntityList{} }
func (m *EntityList) String() string { return proto.CompactTextString(m) }
func (*EntityList) ProtoMessage()    {}
func (*EntityList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae54b1244e5266d5, []int{2}
}

func (m *EntityList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntityList.Unmarshal(m, b)
}
func (m *EntityList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntityList.Marshal(b, m, deterministic)
}
func (m *EntityList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntityList.Merge(m, src)
}
func (m *EntityList) XXX_Size() int {
	return xxx_messageInfo_EntityList.Size(m)
}
func (m *EntityList) XXX_DiscardUnknown() {
	xxx_messageInfo_EntityList.DiscardUnknown(m)
}

var xxx_messageInfo_EntityList proto.InternalMessageInfo

func (m *EntityList) GetOwner() uint64 {
	if m != nil {
		return m.Owner
	}
	return 0
}

func (m *EntityList) GetList() []*EntityContainer {
	if m != nil {
		return m.List
	}
	return nil
}

// 实体增加
type EntityAdd struct {
	Owner                uint64           `protobuf:"varint,1,opt,name=owner,proto3" json:"owner,omitempty"`
	Container            *EntityContainer `protobuf:"bytes,2,opt,name=container,proto3" json:"container,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *EntityAdd) Reset()         { *m = EntityAdd{} }
func (m *EntityAdd) String() string { return proto.CompactTextString(m) }
func (*EntityAdd) ProtoMessage()    {}
func (*EntityAdd) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae54b1244e5266d5, []int{3}
}

func (m *EntityAdd) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntityAdd.Unmarshal(m, b)
}
func (m *EntityAdd) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntityAdd.Marshal(b, m, deterministic)
}
func (m *EntityAdd) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntityAdd.Merge(m, src)
}
func (m *EntityAdd) XXX_Size() int {
	return xxx_messageInfo_EntityAdd.Size(m)
}
func (m *EntityAdd) XXX_DiscardUnknown() {
	xxx_messageInfo_EntityAdd.DiscardUnknown(m)
}

var xxx_messageInfo_EntityAdd proto.InternalMessageInfo

func (m *EntityAdd) GetOwner() uint64 {
	if m != nil {
		return m.Owner
	}
	return 0
}

func (m *EntityAdd) GetContainer() *EntityContainer {
	if m != nil {
		return m.Container
	}
	return nil
}

// 实体更新
type EntityUpdate struct {
	Owner                uint64           `protobuf:"varint,1,opt,name=owner,proto3" json:"owner,omitempty"`
	Container            *EntityContainer `protobuf:"bytes,2,opt,name=container,proto3" json:"container,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *EntityUpdate) Reset()         { *m = EntityUpdate{} }
func (m *EntityUpdate) String() string { return proto.CompactTextString(m) }
func (*EntityUpdate) ProtoMessage()    {}
func (*EntityUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae54b1244e5266d5, []int{4}
}

func (m *EntityUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntityUpdate.Unmarshal(m, b)
}
func (m *EntityUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntityUpdate.Marshal(b, m, deterministic)
}
func (m *EntityUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntityUpdate.Merge(m, src)
}
func (m *EntityUpdate) XXX_Size() int {
	return xxx_messageInfo_EntityUpdate.Size(m)
}
func (m *EntityUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_EntityUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_EntityUpdate proto.InternalMessageInfo

func (m *EntityUpdate) GetOwner() uint64 {
	if m != nil {
		return m.Owner
	}
	return 0
}

func (m *EntityUpdate) GetContainer() *EntityContainer {
	if m != nil {
		return m.Container
	}
	return nil
}

// 实体删除
type EntityDelete struct {
	Owner                uint64      `protobuf:"varint,1,opt,name=owner,proto3" json:"owner,omitempty"`
	EntityType           entity.TYPE `protobuf:"varint,2,opt,name=entityType,proto3,enum=entity.TYPE" json:"entityType,omitempty"`
	List                 []uint64    `protobuf:"varint,3,rep,packed,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *EntityDelete) Reset()         { *m = EntityDelete{} }
func (m *EntityDelete) String() string { return proto.CompactTextString(m) }
func (*EntityDelete) ProtoMessage()    {}
func (*EntityDelete) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae54b1244e5266d5, []int{5}
}

func (m *EntityDelete) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntityDelete.Unmarshal(m, b)
}
func (m *EntityDelete) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntityDelete.Marshal(b, m, deterministic)
}
func (m *EntityDelete) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntityDelete.Merge(m, src)
}
func (m *EntityDelete) XXX_Size() int {
	return xxx_messageInfo_EntityDelete.Size(m)
}
func (m *EntityDelete) XXX_DiscardUnknown() {
	xxx_messageInfo_EntityDelete.DiscardUnknown(m)
}

var xxx_messageInfo_EntityDelete proto.InternalMessageInfo

func (m *EntityDelete) GetOwner() uint64 {
	if m != nil {
		return m.Owner
	}
	return 0
}

func (m *EntityDelete) GetEntityType() entity.TYPE {
	if m != nil {
		return m.EntityType
	}
	return entity.TYPE_None
}

func (m *EntityDelete) GetList() []uint64 {
	if m != nil {
		return m.List
	}
	return nil
}

// 通用错误消息
type CommonError struct {
	OpCode               int32    `protobuf:"varint,1,opt,name=opCode,proto3" json:"opCode,omitempty"`
	ErrCode              int32    `protobuf:"varint,2,opt,name=errCode,proto3" json:"errCode,omitempty"`
	ErrText              string   `protobuf:"bytes,3,opt,name=errText,proto3" json:"errText,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommonError) Reset()         { *m = CommonError{} }
func (m *CommonError) String() string { return proto.CompactTextString(m) }
func (*CommonError) ProtoMessage()    {}
func (*CommonError) Descriptor() ([]byte, []int) {
	return fileDescriptor_ae54b1244e5266d5, []int{6}
}

func (m *CommonError) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonError.Unmarshal(m, b)
}
func (m *CommonError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonError.Marshal(b, m, deterministic)
}
func (m *CommonError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonError.Merge(m, src)
}
func (m *CommonError) XXX_Size() int {
	return xxx_messageInfo_CommonError.Size(m)
}
func (m *CommonError) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonError.DiscardUnknown(m)
}

var xxx_messageInfo_CommonError proto.InternalMessageInfo

func (m *CommonError) GetOpCode() int32 {
	if m != nil {
		return m.OpCode
	}
	return 0
}

func (m *CommonError) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *CommonError) GetErrText() string {
	if m != nil {
		return m.ErrText
	}
	return ""
}

func init() {
	proto.RegisterType((*Ping)(nil), "client.Ping")
	proto.RegisterType((*EntityContainer)(nil), "client.EntityContainer")
	proto.RegisterType((*EntityList)(nil), "client.EntityList")
	proto.RegisterType((*EntityAdd)(nil), "client.EntityAdd")
	proto.RegisterType((*EntityUpdate)(nil), "client.EntityUpdate")
	proto.RegisterType((*EntityDelete)(nil), "client.EntityDelete")
	proto.RegisterType((*CommonError)(nil), "client.CommonError")
}

func init() { proto.RegisterFile("proto/client/common.proto", fileDescriptor_ae54b1244e5266d5) }

var fileDescriptor_ae54b1244e5266d5 = []byte{
	// 389 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x51, 0x8b, 0xda, 0x40,
	0x14, 0x85, 0x31, 0x89, 0xb1, 0xb9, 0x49, 0x2d, 0x0c, 0xa5, 0xa6, 0xed, 0x43, 0x25, 0xf4, 0x41,
	0x68, 0x89, 0x60, 0xe9, 0x0f, 0x68, 0xad, 0x6f, 0x85, 0xca, 0x68, 0xa1, 0xd2, 0x07, 0x89, 0xc9,
	0x6d, 0x08, 0xea, 0x4c, 0x98, 0x0c, 0xe9, 0xfa, 0x63, 0xf7, 0xbf, 0x2c, 0xde, 0x49, 0x56, 0x5d,
	0x70, 0xd9, 0x87, 0x7d, 0xcb, 0xbd, 0xdf, 0x39, 0x87, 0xcc, 0xc9, 0x04, 0xde, 0x96, 0x4a, 0x6a,
	0x39, 0x4e, 0x77, 0x05, 0x0a, 0x3d, 0x4e, 0xe5, 0x7e, 0x2f, 0x45, 0x4c, 0x3b, 0xe6, 0x9a, 0xe5,
	0xbb, 0x46, 0x82, 0x42, 0x17, 0xfa, 0x70, 0x21, 0x89, 0x6e, 0x3b, 0xe0, 0xcc, 0x0b, 0x91, 0xb3,
	0x0f, 0xe0, 0x1b, 0xf5, 0xba, 0x42, 0x91, 0x85, 0x9d, 0x61, 0x67, 0xe4, 0x72, 0x30, 0xab, 0x05,
	0x8a, 0x8c, 0xbd, 0x07, 0x2f, 0x4f, 0x34, 0xae, 0x15, 0xa6, 0x75, 0x68, 0x11, 0x7e, 0x71, 0x5c,
	0x70, 0x4c, 0xeb, 0x7b, 0x48, 0x5e, 0xfb, 0x04, 0xc9, 0x39, 0x80, 0x5e, 0x5e, 0x19, 0x9f, 0x43,
	0xc8, 0xcd, 0x2b, 0x72, 0x19, 0x40, 0x9e, 0x6e, 0x0b, 0xc8, 0x11, 0xc1, 0x4b, 0x8a, 0xdb, 0x24,
	0xe9, 0x76, 0xad, 0xb0, 0x0e, 0x5d, 0xc2, 0xfe, 0x71, 0xf9, 0x3d, 0x49, 0xb7, 0x1c, 0x6b, 0xf6,
	0x11, 0xfa, 0x27, 0x0d, 0x65, 0xf4, 0x48, 0x14, 0xb4, 0xa2, 0x63, 0x52, 0xb4, 0x80, 0x57, 0x33,
	0x3a, 0xf6, 0x54, 0x0a, 0x9d, 0x14, 0x02, 0x15, 0xfb, 0x0c, 0x60, 0x9a, 0x58, 0x1e, 0x4a, 0xa4,
	0x83, 0xf6, 0x27, 0x41, 0x6c, 0x56, 0xf1, 0x72, 0x35, 0x9f, 0xf1, 0x33, 0xce, 0x18, 0x38, 0xbb,
	0xa2, 0xd2, 0xa1, 0x35, 0xb4, 0x47, 0x01, 0xa7, 0xe7, 0xe8, 0x17, 0x80, 0x09, 0xfd, 0x59, 0x54,
	0x9a, 0xbd, 0x86, 0xae, 0xfc, 0x2f, 0x50, 0x51, 0x94, 0xc3, 0xcd, 0xc0, 0x3e, 0x9d, 0xf9, 0xfc,
	0xc9, 0x20, 0x36, 0x4d, 0xc6, 0x0f, 0x5e, 0xa6, 0x09, 0xfc, 0x03, 0x9e, 0x01, 0xdf, 0xb2, 0xec,
	0x4a, 0xde, 0x57, 0xf0, 0xd2, 0xd6, 0x45, 0xf5, 0x3f, 0x12, 0x7a, 0x52, 0x46, 0x7f, 0x21, 0x30,
	0xf4, 0x77, 0x99, 0x25, 0x1a, 0x9f, 0x37, 0xfc, 0x5f, 0x1b, 0xfe, 0x03, 0x77, 0x78, 0x35, 0xfc,
	0xb2, 0x6f, 0xeb, 0x89, 0x7d, 0xdb, 0x43, 0x7b, 0xe4, 0x34, 0xf5, 0xac, 0xc0, 0x9f, 0xd2, 0xa5,
	0x9d, 0x29, 0x25, 0x15, 0x7b, 0x03, 0xae, 0x2c, 0xa7, 0x32, 0x33, 0x1f, 0xaf, 0xcb, 0x9b, 0x89,
	0x85, 0xd0, 0x43, 0xa5, 0x08, 0x58, 0x04, 0xda, 0xb1, 0x21, 0x4b, 0xbc, 0xd1, 0x74, 0x39, 0x3d,
	0xde, 0x8e, 0x1b, 0x97, 0x7e, 0x83, 0x2f, 0x77, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa2, 0x94, 0x18,
	0x39, 0x46, 0x03, 0x00, 0x00,
}
