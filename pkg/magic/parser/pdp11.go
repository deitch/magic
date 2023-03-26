package parser

import "encoding/binary"

// implement ByteOrder for the weird pdp11 middle-endian format, where 16-bit (2-byte) words are
// little-endian, but anything bigger is combinations of the 2-byte words in big-endian order.

var Pdp11ByteOrder pdp11ByteOrder

type pdp11ByteOrder struct {
}

func (p pdp11ByteOrder) Uint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)

}
func (p pdp11ByteOrder) Uint32(b []byte) uint32 {
	return uint32(binary.LittleEndian.Uint16(b[2:4])) | uint32(binary.LittleEndian.Uint16(b[0:2]))<<16
}
func (p pdp11ByteOrder) Uint64(b []byte) uint64 {
	return uint64(p.Uint16(b[6:8])) |
		uint64(p.Uint16(b[4:6]))<<16 |
		uint64(p.Uint16(b[2:4]))<<32 |
		uint64(p.Uint16(b[0:2]))<<48
}
func (p pdp11ByteOrder) PutUint16(b []byte, u uint16) {
	binary.LittleEndian.PutUint16(b, u)
}

func (p pdp11ByteOrder) PutUint32(b []byte, u uint32) {
	p.PutUint16(b[0:2], uint16(u>>16))
	p.PutUint16(b[2:4], uint16(u))
}
func (p pdp11ByteOrder) PutUint64(b []byte, u uint64) {
	p.PutUint16(b[0:2], uint16(u>>48))
	p.PutUint16(b[2:4], uint16(u>>32))
	p.PutUint16(b[4:6], uint16(u>>16))
	p.PutUint16(b[6:8], uint16(u))
}
func (p pdp11ByteOrder) String() string {
	return "Pdp11ByteOrder"
}
