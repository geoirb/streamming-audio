package wav

import (
	"bytes"
	"io"

	"github.com/cryptix/wav"
)

// WAV audio file
type WAV struct{}

// Audio parse wav file
func (w *WAV) Audio(data []byte) (reader io.Reader, channels uint16, rate uint32, err error) {
	f := bytes.NewReader(data)
	wavReader, err := wav.NewReader(f, f.Size())
	if err != nil {
		return
	}
	channels = wavReader.GetNumChannels()
	rate = wavReader.GetSampleRate()
	reader, err = wavReader.GetDumbReader()
	return
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
