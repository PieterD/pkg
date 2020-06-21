package shorthand

import (
	"encoding/binary"
	"hash/crc32"
	"sync"
)

type Encoder struct {
	b      []byte
	crcPos int
}

var bytesPool = sync.Pool{New: func() interface{} { return make([]byte, 1024) }}

func NewEncoder() *Encoder {
	return &Encoder{
		b: bytesPool.Get().([]byte)[:0],
	}
}

func (e *Encoder) Close() {
	bytesPool.Put(e.b)
}

func (e *Encoder) Advance(size int) {
	e.b = e.b[:len(e.b)+size]
}

func (e *Encoder) Grow(required int) []byte {
	if cap(e.b) < len(e.b)+required {
		e.b = append(e.b, make([]byte, required)...)
	}
	return e.b[len(e.b) : len(e.b)+required]
}

func (e *Encoder) Copy() []byte {
	c := make([]byte, len(e.b), len(e.b))
	copy(c, e.b)
	return c
}

func (e *Encoder) Buffer() []byte {
	return e.b
}

func (e *Encoder) Uint8(i uint8) {
	e.b = append(e.b, i)
}

func (e *Encoder) Uint16(i uint16) {
	binary.BigEndian.PutUint16(e.Grow(2), i)
	e.Advance(2)
}

func (e *Encoder) Uint32(i uint32) {
	binary.BigEndian.PutUint32(e.Grow(4), i)
	e.Advance(4)
}

func (e *Encoder) Uint64(i uint64) {
	binary.BigEndian.PutUint64(e.Grow(8), i)
	e.Advance(8)
}

func (e *Encoder) VarInt(i int) {
	e.VarInt64(int64(i))
}

func (e *Encoder) VarInt64(i int64) {
	e.Advance(binary.PutVarint(e.Grow(binary.MaxVarintLen64), i))
}

func (e *Encoder) VarUint64(u uint64) {
	e.Advance(binary.PutUvarint(e.Grow(binary.MaxVarintLen64), u))
}

func (e *Encoder) Bytes(b []byte) {
	e.b = append(e.b, b...)
}

func (e *Encoder) ByteSlice(b []byte) {
	e.VarInt(len(b))
	e.Bytes(b)
}

func (e *Encoder) String(s string) {
	e.ByteSlice([]byte(s))
}

func (e *Encoder) StartCRC() {
	e.crcPos = len(e.b)
}

func (e *Encoder) PutCRC() {
	e.Uint32(crc32.ChecksumIEEE(e.b[e.crcPos:]))
}
