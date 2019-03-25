package net

import "testing"

func TestByteBuffer(t *testing.T) {
	var b = NewByteBuffer(2)
	b.AppendUint8(0xE)
	b.AppendUint16(0xED)
	b.AppendUint32(0xEDAB)
	b.AppendUint64(0xABCDEF12)
	b.AppendByte('a')
	b.AppendBytes([]byte("abc"))
	b.AppendString("efg")
	b.AppendVarint(10)
	b.AppendVarint(777)
	b.AppendVarint(1000002882992)

	if v, _ := b.ReadUint8(); v != 0xE {
		t.Fatal(1, b)
	}
	if v, _ := b.ReadUint16(); v != 0xED {
		t.Fatal(2)
	}
	if v, _ := b.ReadUint32(); v != 0xEDAB {
		t.Fatal(3)
	}
	if v, _ := b.ReadUint64(); v != 0xABCDEF12 {
		t.Fatal(4)
	}
	if v, _ := b.ReadByte(); v != 'a' {
		t.Fatal(5)
	}
	if v, _ := b.ReadBytes(3); string(v) != "abc" {
		t.Fatal(v)
	}
	if v, _ := b.ReadString(3); v != "efg" {
		t.Fatal(7)
	}
	if v, _ := b.ReadVarintUint16(); v != 10 {
		t.Fatal(8)
	}
	if v, _ := b.ReadVarintUint32(); v != 777 {
		t.Fatal(9)
	}
	if v, _ := b.ReadVarintUint64(); v != 1000002882992 {
		t.Fatal(10)
	}
}
