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
func (w *WAV) Write(ctx context.Context, name string, channels uint16, rate uint32) (io.Writer, error) {
	if !strings.HasSuffix(name, ".wav") {
		name = name + ".wav"
	}
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	wav := wav.NewWriter(file, channels, rate, wav.FormatS16LE)
	go func() {
		<-ctx.Done()
		wav.Close()
	}()

	return wav, nil
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
