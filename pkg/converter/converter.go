package converter

import "encoding/binary"

// Converter ...
type Converter struct {
}

// ToInt16 convert array byte to array int16
func (c *Converter) ToInt16(src []byte) (dst []int16) {
	dst = make([]int16, 0, len(src)/2)
	for i := 0; i+2 <= len(src); i += 2 {
		buf := src[i : i+2]
		dst = append(dst, int16(binary.LittleEndian.Uint16(buf)))
	}
	return
}

// NewConverter ...
func NewConverter() *Converter {
	return &Converter{}
}
