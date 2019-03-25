package net

import (
	"core/log"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	protoHelp "proto"
)

const MAX_PACKET_LENGTH = 32 * 1024

// 1- 99是保留op code
const (
	OP_C2G_HEART_BEAT = 1 // Client到Gate的心跳
	OP_G2C_HEART_BEAT = 2 // Gate到Client心跳
	OP_G2S_HEART_BEAT = 3 // Gate到后端Server心跳
	OP_S2G_HEART_BEAT = 4 // 后端Server到Gate心跳
)

type IMsgHeader interface {
	GetLength() int
	GetOpCode() uint16
	PeekSize(b *ByteBuffer)
	Decode(b *ByteBuffer)
	Encode(b *ByteBuffer)
}

func NewInternalHeader(opCode uint16, guid uint64) *InternalHeader {
	return &InternalHeader{OpCode: opCode, Guid: guid}
}

// 服务器内部通讯接口 gate<->server
type InternalHeader struct {
	Size    uint32 // 消息体大小，包含包头
	OpCode  uint16 // 消息ID
	Guid    uint64 // guid
	OriSize uint32 // 数据体未压缩大小
}

func (i *InternalHeader) GetLength() int {
	return 4 + 2 + 8 + 4
}
func (i *InternalHeader) GetOpCode() uint16 {
	return i.OpCode
}
func (i *InternalHeader) PeekSize(b *ByteBuffer) {
	i.Size, _ = b.PeekUint32()
}

func (i *InternalHeader) Decode(b *ByteBuffer) {
	i.Size, _ = b.ReadUint32()
	i.OpCode, _ = b.ReadUint16()
	i.Guid, _ = b.ReadUint64()
	i.OriSize, _ = b.ReadUint32()
}
func (i *InternalHeader) Encode(b *ByteBuffer) {
	b.AppendUint32(i.Size)
	b.AppendUint16(i.OpCode)
	b.AppendUint64(i.Guid)
	b.AppendUint32(i.OriSize)
}

func NewClientHeaderIn(opCode uint16, guid uint64) *ClientHeaderIn {
	return &ClientHeaderIn{OpCode: opCode}
}

// 接收到的客户端消息头  client->gate
type ClientHeaderIn struct {
	Size     uint32 // 消息体大小，包含包头
	OpCode   uint16 // 消息ID
	CheckSum uint8  // checksum
}

func (c *ClientHeaderIn) GetLength() int {
	return 4 + 2 + 1
}
func (c *ClientHeaderIn) GetOpCode() uint16 {
	return c.OpCode
}

func (c *ClientHeaderIn) PeekSize(b *ByteBuffer) {
	c.Size, _ = b.PeekUint32()
}

func (c *ClientHeaderIn) Decode(b *ByteBuffer) {
	c.Size, _ = b.ReadUint32()
	c.OpCode, _ = b.ReadUint16()
	c.CheckSum, _ = b.ReadUint8()

}
func (c *ClientHeaderIn) Encode(b *ByteBuffer) {
	b.AppendUint32(c.Size)
	b.AppendUint16(c.OpCode)
	b.AppendUint8(c.CheckSum)
}

func NewClientHeaderOut(opCode uint16, guid uint64) *ClientHeaderOut {
	return &ClientHeaderOut{OpCode: opCode}
}

// 发送到客户端的消息头  gate->client
type ClientHeaderOut struct {
	Size     uint32 // 消息体大小，包含包头
	OpCode   uint16 // 消息ID
	OriSize  uint32 // 数据体未压缩大小
	CheckSum uint8  // checksum
}

func (c *ClientHeaderOut) GetLength() int {
	return 4 + 2 + 4 + 1
}
func (c *ClientHeaderOut) GetOpCode() uint16 {
	return c.OpCode
}

func (c *ClientHeaderOut) PeekSize(b *ByteBuffer) {
	c.Size, _ = b.PeekUint32()
}

func (c *ClientHeaderOut) Decode(b *ByteBuffer) {
	c.Size, _ = b.ReadUint32()
	c.OpCode, _ = b.ReadUint16()
	c.OriSize, _ = b.ReadUint32()
	c.CheckSum, _ = b.ReadUint8()
}
func (c *ClientHeaderOut) Encode(b *ByteBuffer) {
	b.AppendUint32(c.Size)
	b.AppendUint16(c.OpCode)
	b.AppendUint32(c.OriSize)
	b.AppendUint8(c.CheckSum)
}

func NewPacket(opCode uint16, msgBody interface{}, guid uint64) *Packet {
	head := &InternalHeader{OpCode: opCode, Guid: guid}
	return &Packet{IMsgHeader: head, MsgBody: msgBody}
}
func NewPacketWithData(opCode uint16, b *ByteBuffer, guid uint64) *Packet {
	head := &InternalHeader{OpCode: opCode, Guid: guid}
	return &Packet{IMsgHeader: head, Buff: b}
}

func NewPacketToClient(opCode uint16, msgBody interface{}) *Packet {
	head := &ClientHeaderOut{OpCode: opCode}
	return &Packet{IMsgHeader: head, MsgBody: msgBody}
}

func NewPackeToClienttWithData(opCode uint16, b *ByteBuffer) *Packet {
	head := &ClientHeaderOut{OpCode: opCode}
	return &Packet{IMsgHeader: head, Buff: b}
}

//    以下是经过io线程处理后的packet对象的结构，包含收到和发送。
//    ---------------------Gate收到的packet --------------------
//    |   IMsgHeader    |   MsgBody   |   Buff                 |
//    | ClientHeaderIn  |   nil       |   pb结构对应的二进制     |
//    ---------------------Gate发出的packet --------------------
//    |   IMsgHeader    |   MsgBody   |    Buff                |
//    | ClientHeaderOut |   nil       |    head + pb对应的二进制 |
//    --------内部收到的packet（自动反序列化，需要在pb_help注册）----
//    |   IMsgHeader    |   MsgBody   |    Buff                |
//    | InternalHeader  |   pb对象     |    nil                 |
//    ---------------------内部发出的packet-----------------------
//    |   IMsgHeader    |   MsgBody   |   Buff                  |
//    | InternalHeader  |   nil       |   head+pb 对应的二进制   |
//    TODO： 广播包是否需要单独优化，把子包的序列化也放在io线程
type Packet struct {
	IMsgHeader             // 逻辑操作的包头结构
	MsgBody    interface{} // 逻辑操作的数据体,比如 protobuf 对象
	Buff       *ByteBuffer // io发送和读取操作的对象，二进制
}

// Client -> Gate
// Gate收到客户端消息时，在io线程的预处理逻辑,包头为 ClientHeaderIn
func PreReadClientToGateHook(b *ByteBuffer) (*Packet, error) {
	head := &ClientHeaderIn{}
	buffSize := b.Length()
	if buffSize < head.GetLength() {
		return nil, nil
	}
	head.PeekSize(b)
	if head.Size > MAX_PACKET_LENGTH {
		return nil, errors.New(fmt.Sprintf("packet length is to large:%d", head.Size))
	}
	if uint32(buffSize) < head.Size {
		return nil, nil
	}
	head.Decode(b)
	pkt := &Packet{IMsgHeader: head}
	// 客户端发来的消息，不进行消息体反序列化，因为绝大部分消息都是转发。
	bodySize := int(head.Size) - head.GetLength()
	if bodySize > 0 {
		pkt.Buff = b.Copy(bodySize) // 直接从copy到buffer中返回
	}
	return pkt, nil
}

//Gate -> Client
//Gate 发送给client,在io线程的预处理逻辑，包头为 ClientHeaderOut
func PreWriteGateToClientHook(p *Packet) {
	head := p.IMsgHeader.(*ClientHeaderOut)
	if p.MsgBody != nil { // Gate自己构建的包
		buff, err := proto.Marshal(p.MsgBody.(proto.Message)) //TODO：优化，直接在ByteBuffer上Marshal 减少一次copy
		if err != nil {
			log.Error("[PACKET]:marshal message error", log.Uint16("OpCode", head.GetOpCode()))
		}
		head.Size = uint32(head.GetLength() + len(buff))
		p.Buff = NewByteBuffer(int(head.Size))
		head.Encode(p.Buff)
		p.Buff.Write(buff)
	} else { // Gate转发的包 或者是空包
		if p.Buff != nil {
			size := head.GetLength() + p.Buff.Length()
			head.Size = uint32(size)
			buf := NewByteBuffer(int(head.Size))
			head.Encode(buf)
			buf.Concat(p.Buff)
			p.Buff = buf
		} else { // 空包
			head.Size = uint32(head.GetLength())
			p.Buff = NewByteBuffer(int(head.Size))
			head.Encode(p.Buff)
		}
	}
}

//Gate<->Server
//Gate 和后端交互时，在io线程的预处理逻辑，包头为 InternalHeader
func PreReadInternalHook(b *ByteBuffer) (*Packet, error) {
	head := &InternalHeader{}
	buffSize := b.Length()
	if buffSize < head.GetLength() {
		return nil, nil
	}
	head.PeekSize(b)
	if head.Size > MAX_PACKET_LENGTH {
		return nil, errors.New(fmt.Sprintf("packet length is to large:%d", head.Size))
	}
	if uint32(buffSize) < head.Size {
		return nil, nil
	}
	head.Decode(b)
	pkt := &Packet{IMsgHeader: head}

	bodySize := int(head.Size) - head.GetLength()
	if bodySize > 0 { // 如果携带了消息体，继续反序列化消息体
		newFunc := protoHelp.OpCodeToMessage[head.OpCode]
		if newFunc != nil {
			msgBody := newFunc()
			err := proto.Unmarshal(b.Read(bodySize), msgBody)
			if err != nil {
				return nil, err
			}
			pkt.MsgBody = msgBody
			pkt.Buff = nil
		} else { // 找不到自动反序列化函数
			log.Warn("[PACKET]:auto unmarshal function no found", log.Uint16("OpCode", head.OpCode))
			pkt.Buff = b.Copy(bodySize)
		}
	}
	return pkt, nil
}

//Gate<->Server
//Gate 和后端交互时，在io线程的预处理逻辑，包头为 InternalHeader
func PreWriteInternalHook(p *Packet) {
	head := p.IMsgHeader.(*InternalHeader)
	if p.MsgBody != nil { // 自己构建的包
		buff, err := proto.Marshal(p.MsgBody.(proto.Message)) //TODO：优化，直接在ByteBuffer上Marshal 减少一次copy
		if err != nil {
			log.Error("[PACKET]:marshal message error", log.Uint16("OpCode", head.GetOpCode()))
		}
		head.Size = uint32(head.GetLength() + len(buff))
		p.Buff = NewByteBuffer(int(head.Size))
		head.Encode(p.Buff)
		p.Buff.Write(buff)
	} else { // 转发的包 或者是空包
		if p.Buff != nil {
			size := head.GetLength() + p.Buff.Length()
			head.Size = uint32(size)
			buf := NewByteBuffer(int(head.Size))
			head.Encode(buf)
			buf.Concat(p.Buff)
			p.Buff = buf
		} else { // 空包
			head.Size = uint32(head.GetLength())
			p.Buff = NewByteBuffer(int(head.Size))
			head.Encode(p.Buff)
		}
	}
}
