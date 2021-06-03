package converter

import (
	"encoding/binary"
)

const bytePerInt16 = 2

// Converter ...
type Converter struct{}

// ToInt16 convert array byte to array int16
func (c *Converter) ToInt16(src []byte) (dst []int16) {
	dst = make([]int16, 0, len(src)/bytePerInt16)
	for i := 0; i <= len(src)-bytePerInt16; i += bytePerInt16 {
		dst = append(
			dst,
			int16(binary.LittleEndian.Uint16(src[i:i+bytePerInt16])),
		)
	}
	return
}

// ToByte convert array int16 to array byte
func (c *Converter) ToByte(src []int16) (dst []byte) {
	dst = make([]byte, 0, len(src)*bytePerInt16)
	b := make([]byte, bytePerInt16)
	for _, i := range src {
		binary.LittleEndian.PutUint16(b, uint16(i))
		dst = append(dst, b...)
	}
	return
}

// NewConverter ...
func NewConverter() *Converter {
	return &Converter{}
}
