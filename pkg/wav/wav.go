package wav

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/cryptix/wav"
)

// WAV audio file
type WAV struct {
	bufferSize int
}

// Read wav file
func (w *WAV) Read(data []byte) (reader io.Reader, channels uint16, rate uint32, err error) {
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

// Recode wav file
func (w *WAV) Recode(ctx context.Context, name string, channels uint16, rate uint32, r io.Reader) (err error) {
	if !strings.HasSuffix(name, ".wav") {
		name = name + ".wav"
	}
	file, err := os.Create(name)
	if err != nil {
		return
	}
	defer file.Close()

	meta := wav.File{
		Channels:        channels,
		SampleRate:      rate,
		SignificantBits: 16,
	}

	writer, err := meta.NewWriter(file)
	if err != nil {
		return
	}
	defer writer.Close()

	buffer := make([]byte, w.bufferSize)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if l, err := r.Read(buffer); err == nil {
					if _, err = writer.Write(buffer[:l]); err != nil {
						return
					}
				}
			}
		}
	}()
	return
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
