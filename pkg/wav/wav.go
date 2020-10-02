package wav

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/geoirb/wav"
)

// WAV audio file
type WAV struct{}

// Read wav file
func (w *WAV) Read(data []byte) (reader io.Reader, channels uint16, rate uint32, err error) {
	wav, err := wav.NewReader(data)
	if err != nil {
		return
	}
	channels = wav.GetNumChannels()
	rate = wav.GetSampleRate()
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

	writer = wav.NewWriter(file)
	go func() {
		<-ctx.Done()
		audio.Close()
	}()

	return
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
