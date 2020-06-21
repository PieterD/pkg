package shorthand

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"github.com/PieterD/pkg/panic"
)

const (
	maxUint = uint64(^uint(0))
	maxInt  = int64(int(maxUint >> 1))
	minInt  = int64(-maxInt - 1)
)

var (
	ErrInvalidCRC = fmt.Errorf("invalid CRC")
)

func Recover(errp *error) {
	if err := panic.RecoverCheck(recover()); err != nil {
		*errp = err
	}
}

type Decoder struct {
	buf    []byte
	pos    int
	crcPos int
}

func NewDecoder(b []byte) *Decoder {
	return &Decoder{
		buf: b,
	}
}

func (d *Decoder) Len() int {
	return len(d.buf) - d.pos
}

func (d *Decoder) Buffer() []byte {
	return d.buf[d.pos:]
}

func (d *Decoder) Available(fun, field string, num int) []byte {
	if d.pos+num > len(d.buf) {
		panic.Panic(fmt.Errorf("%s(%s) not enough bytes, needed %d", fun, field, num))
	}
	return d.buf[d.pos : d.pos+num]
}

func (d *Decoder) Advance(fun, field string, n int) {
	d.Available(fun, field, n)
	d.pos += n
}

func (d *Decoder) Uint8(field string) uint8 {
	b := d.Available("Uint8", field, 1)
	i := b[0]
	d.Advance("Uint8", field, 1)
	return i
}

func (d *Decoder) Uint16(field string) uint16 {
	b := d.Available("Uint16", field, 2)
	i := binary.BigEndian.Uint16(b)
	d.Advance("Uint16", field, 2)
	return i
}

func (d *Decoder) Uint32(field string) uint32 {
	b := d.Available("Uint32", field, 4)
	i := binary.BigEndian.Uint32(b)
	d.Advance("Uint32", field, 4)
	return i
}

func (d *Decoder) Uint64(field string) uint64 {
	b := d.Available("Uint64", field, 8)
	i := binary.BigEndian.Uint64(b)
	d.Advance("Uint64", field, 8)
	return i
}

func (d *Decoder) VarInt(field string) int {
	i := d.VarInt64(field)
	if i > maxInt {
		panic.Panic(fmt.Errorf("VarInt(%s) %d too large (>%d) for int", field, i, maxInt))
	} else if i < minInt {
		panic.Panic(fmt.Errorf("VarInt(%s) %d too small (<%d) for int", field, i, minInt))
	}
	return int(i)
}

func (d *Decoder) VarInt64(field string) int64 {
	i, n := binary.Varint(d.Buffer())
	if n == 0 {
		panic.Panic(fmt.Errorf("VarInt64(%s) buffer too small", field))
	}
	if n < 0 {
		panic.Panic(fmt.Errorf("VarInt64(%s) overflow(%d)", field, n))
	}
	d.Advance("VarInt64", field, n)
	return i
}

func (d *Decoder) VarUint64(field string) uint64 {
	i, n := binary.Uvarint(d.Buffer())
	if n == 0 {
		panic.Panic(fmt.Errorf("VarInt64(%s) buffer too small", field))
	}
	if n < 0 {
		panic.Panic(fmt.Errorf("VarInt64(%s) overflow(%d)", field, n))
	}
	d.Advance("VarUint64", field, n)
	return i
}

func (d *Decoder) Bytes(field string, num int) []byte {
	b := d.Available("Bytes", field, num)
	cop := make([]byte, num)
	copy(cop, b)
	d.Advance("Bytes", field, num)
	return b
}

func (d *Decoder) ByteSlice(field string) []byte {
	size := d.VarInt(field)
	return d.Bytes(field, size)
}

func (d *Decoder) String(field string) string {
	return string(d.ByteSlice(field))
}

func (d *Decoder) StartCRC() {
	d.crcPos = d.pos
}

func (d *Decoder) CheckCRC(field string) {
	_ = d.Available("CheckCRC", field, 4)
	got := crc32.ChecksumIEEE(d.buf[d.crcPos:d.pos])
	want := d.Uint32(field)
	if want != got {
		panic.Panic(fmt.Errorf("invalid crc on field %s: %w", field, ErrInvalidCRC))
	}
}
