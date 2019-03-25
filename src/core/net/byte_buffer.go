package net

import "io"

func NewByteBuffer(size int) *ByteBuffer {
	var b ByteBuffer
	b.malloc(size)
	return &b
}

type ByteBuffer struct {
	buf []byte
	r   int
	w   int
}

func (b *ByteBuffer) malloc(size int) {
	if b.buf != nil {
		if size > len(b.buf) {
			if size <= cap(b.buf) { //cap够用的时候，避免重新allocate
				b.buf = b.buf[:size]
			} else {
				buf := make([]byte, size, size*2)
				copy(buf, b.buf)
				b.buf = buf
			}
		}
	} else if size > 0 {
		b.buf = make([]byte, size)
		b.Reset()
	}
}

func (b *ByteBuffer) Realloc(size int) {
	b.malloc(size)
}

func (b *ByteBuffer) Read(size int) (x []byte) {
	x = b.buf[b.r : b.r+size]
	b.r += size
	return
}

func (b *ByteBuffer) Write(v []byte) {
	var n = copy(b.buf[b.w:], v)
	b.w += n
}

//返回当前可读的字节
func (b *ByteBuffer) Data() []byte {
	return b.buf[b.r:b.w]
}

//当前可读字节数
func (b *ByteBuffer) Length() int {
	return b.w - b.r
}

//内部缓冲区长度
func (b *ByteBuffer) Size() int {
	return len(b.buf)
}

//剩余容量
func (b *ByteBuffer) Space() int {
	return len(b.buf) - b.w
}

//删除已读的数据
func (b *ByteBuffer) Slim() {
	copy(b.buf, b.buf[b.r:b.w])
	b.w = b.w - b.r
	b.r = 0
}

//连接buffer
func (b *ByteBuffer) Concat(o *ByteBuffer) {
	n := copy(b.buf[b.w:], o.Data())
	b.w += n
}

func (b *ByteBuffer) Clone() *ByteBuffer {
	c := NewByteBuffer(b.Size())
	c.w = b.w
	c.r = b.r
	copy(c.buf, b.buf)
	return c
}

// 从buffer中读取一个新copy
func (b *ByteBuffer) Copy(size int) *ByteBuffer {
	c := NewByteBuffer(size)
	n := copy(c.buf, b.buf[b.r:])
	c.w += n
	b.r += n
	return c
}

func (b *ByteBuffer) Reset() {
	b.w = 0
	b.r = 0
}

//使用reader读一次
func (b *ByteBuffer) AppendReader(r io.Reader) error {
	n, err := r.Read(b.buf[b.w:])
	if err != nil {
		return err
	}
	b.w += n
	return nil
}

//输入到一个writer中
func (b *ByteBuffer) AppendWriter(w io.Writer) error {
	if b.Length() <= 0 {
		return nil
	}
	n, err := w.Write(b.buf[b.r:b.w])
	if err != nil {
		return err
	}
	b.r += n
	return nil
}

func (b *ByteBuffer) AppendUint8(v uint8) {
	AppendUint8(b.buf[:b.w], v)
	b.w += 1
}

func (b *ByteBuffer) AppendUint16(v uint16) {
	AppendUint16(b.buf[:b.w], v)
	b.w += 2
}

func (b *ByteBuffer) AppendUint32(v uint32) {
	AppendUint32(b.buf[:b.w], v)
	b.w += 4
}

func (b *ByteBuffer) AppendUint64(v uint64) {
	AppendUint64(b.buf[:b.w], v)
	b.w += 8
}

func (b *ByteBuffer) AppendVarint(v uint64) {
	_, size := EncodeVarint(b.buf[:b.w], v)
	b.w += size
}

func (b *ByteBuffer) AppendString(v string) {
	b.AppendBytes([]byte(v))
}

func (b *ByteBuffer) AppendByte(v byte) {
	b.AppendUint8(v)
}

func (b *ByteBuffer) AppendBytes(v []byte) {
	var n = copy(b.buf[b.w:], v)
	b.w += n
}

func (b *ByteBuffer) ReadUint8() (x uint8, n int) {
	x, n = ReadUint8(b.buf[b.r:b.w])
	b.r += n
	return
}

func (b *ByteBuffer) PeekUint8() (x uint8, n int) {
	x, n = ReadUint8(b.buf[b.r:b.w])
	return
}

func (b *ByteBuffer) ReadUint16() (x uint16, n int) {
	x, n = ReadUint16(b.buf[b.r:b.w])
	b.r += n
	return
}

func (b *ByteBuffer) PeekUint16() (x uint16, n int) {
	x, n = ReadUint16(b.buf[b.r:b.w])
	return
}

func (b *ByteBuffer) ReadUint32() (x uint32, n int) {
	x, n = ReadUint32(b.buf[b.r:b.w])
	b.r += n
	return
}

func (b *ByteBuffer) PeekUint32() (x uint32, n int) {
	x, n = ReadUint32(b.buf[b.r:b.w])
	return
}

func (b *ByteBuffer) ReadUint64() (x uint64, n int) {
	x, n = ReadUint64(b.buf[b.r:b.w])
	b.r += n
	return
}

func (b *ByteBuffer) PeekUint64() (x uint64, n int) {
	x, n = ReadUint64(b.buf[b.r:b.w])
	return
}

func (b *ByteBuffer) ReadVarintUint16() (x uint16, n int) {
	v, n := DecodeVarint(b.buf[b.r:b.w])
	if n <= 3 {
		b.r += n
		return uint16(v), n
	}
	return 0, 0
}

func (b *ByteBuffer) ReadVarintUint32() (x uint32, n int) {
	v, n := DecodeVarint(b.buf[b.r:b.w])
	if n <= 5 {
		b.r += n
		return uint32(v), n
	}
	return 0, 0
}

func (b *ByteBuffer) ReadVarintUint64() (x uint64, n int) {
	v, n := DecodeVarint(b.buf[b.r:b.w])
	if n <= 10 {
		b.r += n
		return uint64(v), n
	}
	return 0, 0
}

func (b *ByteBuffer) ReadByte() (x byte, n int) {
	x, n = ReadUint8(b.buf[b.r:b.w])
	b.r += n
	return
}

func (b *ByteBuffer) PeekByte() (x byte, n int) {
	x, n = ReadUint8(b.buf[b.r:b.w])
	return
}

func (b *ByteBuffer) ReadBytes(size int) (x []byte, n int) {
	x = b.buf[b.r : b.r+size]
	b.r += size
	n = size
	return
}

func (b *ByteBuffer) PeekBytes(size int) (x []byte, n int) {
	x = b.buf[b.r : b.r+size]
	n = size
	return
}

func (b *ByteBuffer) ReadString(size int) (x string, n int) {
	v, n := b.ReadBytes(size)
	return string(v), n
}

func (b *ByteBuffer) PeekString(size int) (x string, n int) {
	v, n := b.PeekBytes(size)
	return string(v), n
}

func AppendUint8(b []byte, v uint8) []byte {
	return append(b, byte(v))
}

func AppendUint16(b []byte, v uint16) []byte {
	return append(b, byte(v>>8), byte(v))
}

func AppendUint32(b []byte, v uint32) []byte {
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func AppendUint64(b []byte, v uint64) []byte {
	b = append(b, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	return b
}

func EncodeVarint(b []byte, v uint64) ([]byte, int) {
	var size = 1
	for v >= 1<<7 {
		b = append(b, byte(v&0x7f|0x80))
		v >>= 7
		size++
	}
	return append(b, byte(v)), size
}

func ReadUint8(b []byte) (x uint8, n int) {
	if len(b) >= 1 {
		return b[0], 1
	}
	return 0, 0
}

func ReadUint16(b []byte) (x uint16, n int) {
	if len(b) >= 2 {
		return uint16(b[0])<<8 | uint16(b[1]), 2
	}
	return 0, 0
}

func ReadUint32(b []byte) (x uint32, n int) {
	if len(b) >= 4 {
		return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) |
			(uint32(b[2]) << 8) | uint32(b[3]), 4
	}
	return 0, 0
}

func ReadUint64(b []byte) (x uint64, n int) {
	if len(b) >= 8 {
		return (uint64(b[0]) << 56) | (uint64(b[1]) << 48) |
			(uint64(b[2]) << 40) | (uint64(b[3]) << 32) |
			(uint64(b[4]) << 24) | (uint64(b[5]) << 16) |
			(uint64(b[6]) << 8) | uint64(b[7]), 8
	}
	return 0, 0
}

func DecodeVarint(b []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(b) {
			return 0, 0
		}
		b := uint64(b[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}
	// The number is too large to represent in a 64-bit value.
	return 0, 0
}
