package wav

import (
	"bytes"
	"io"

	"github.com/cryptix/wav"
)

// WAV audio file
type WAV struct {
	reader         io.Reader
	channels       int
	bytesPerSample int
}

// Parse wav file
func (w *WAV) Parse(data []byte) (err error) {
	f := bytes.NewReader(data)
	wavReader, err := wav.NewReader(f, f.Size())
	if err != nil {
		return
	}
	w.channels = int(wavReader.GetNumChannels())
	w.bytesPerSample = int(wavReader.GetBitsPerSample() / 8)
	w.reader, err = wavReader.GetDumbReader()
	return
}

// Read audio bytes
func (w *WAV) Read() ([]byte, error) {
	//todo
	//min 3
	samples := make([]byte, 1024)
	l, err := w.reader.Read(samples)
	return samples[:l], err
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
