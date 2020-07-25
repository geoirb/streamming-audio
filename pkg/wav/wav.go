package wav

import (
	"bytes"

	"github.com/cryptix/wav"
)

// WAV audio file
type WAV struct {
	reader *wav.Reader
}

// Parse wav file
func (w *WAV) Parse(data []byte) (err error) {
	f := bytes.NewReader(data)
	w.reader, err = wav.NewReader(
		f,
		f.Size(),
	)
	return
}

// GetSample get audio sample
func (w *WAV) GetSample() (samples []int16, err error) {
	var s []int32
	// todo
	// del consts
	s, err = w.reader.ReadSampleEvery(
		2,
		0,
	)

	samples = make([]int16, 0, len(s))
	for _, semple := range s {
		samples = append(samples, int16(semple))
	}

	return
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
