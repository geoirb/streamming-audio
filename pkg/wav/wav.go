package wav

import (
	"bytes"
	"context"
	"io"

	"github.com/cryptix/wav"
)

// WAV audio file
type WAV struct {
	reader         io.Reader
	channels       int
	bytesPerSample int

	sample chan []byte
	err    chan error
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
	w.sample, w.err = make(chan []byte, 1), make(chan error, 1)
	w.reader, err = wavReader.GetDumbReader()
	return
}

// StartReadingSample reading audio sample
func (w *WAV) StartReadingSample(ctx context.Context) {
	sample := make([]byte, w.bytesPerSample*w.channels)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := w.reader.Read(sample); err != nil {
				w.err <- err
				return
			}
			w.sample <- sample
		}
	}
}

// Sample return chan for audio sample
func (w *WAV) Sample() <-chan []byte {
	return w.sample
}

// Error return chan for reading error
func (w *WAV) Error() <-chan error {
	return w.err
}

// StopReadingSample ...
func (w *WAV) StopReadingSample() {
	close(w.sample)
	close(w.err)
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
