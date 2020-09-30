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

// Write wav file
func (w *WAV) Write(ctx context.Context, name string, channels uint16, rate uint32) (writer io.Writer, err error) {
	if !strings.HasSuffix(name, ".wav") {
		name = name + ".wav"
	}
	file, err := os.Create(name)
	if err != nil {
		return
	}

	meta := wav.File{
		Channels:        1,
		SampleRate:      44100,
		SignificantBits: 16,
	}

	audio, err := meta.NewWriter(file)
	if err != nil {
		file.Close()
	}

	go func() {
		<-ctx.Done()
		audio.Close()
	}()
	writer = audio
	return
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
