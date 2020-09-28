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

// Record wav file
func (w *WAV) Record(ctx context.Context, name string, channels uint16, rate uint32, r io.ReadCloser) (err error) {
	if !strings.HasSuffix(name, ".wav") {
		name = name + ".wav"
	}
	file, err := os.Create(name)
	if err != nil {
		r.Close()
		return
	}

	meta := wav.File{
		Channels:        1,
		SampleRate:      rate,
		SignificantBits: 32,
	}

	writer, err := meta.NewWriter(file)
	if err != nil {
		r.Close()
		file.Close()
		return
	}

	go func() {
		<-ctx.Done()
		r.Close()
		file.Close()
		writer.Close()
	}()

	go func() {
		buffer := make([]byte, w.bufferSize)
		for {
			l, err := r.Read(buffer)
			if err != nil {
				return
			}
			if _, err = writer.Write(buffer[:l]); err != nil {
				return
			}
		}
	}()
	return
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
