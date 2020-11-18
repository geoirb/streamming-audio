package wav

import (
	"io"
	"os"
	"strings"

	"github.com/geoirb/wav"
)

// WAV audio file
type WAV struct{}

// Reader wav file
func (w *WAV) Reader(data []byte) (r io.Reader, channels uint16, rate uint32, bitsPerSample uint16, err error) {
	wav, err := wav.NewReader(data)
	if err != nil {
		return
	}
	channels = wav.GetNumChannels()
	rate = wav.GetSampleRate()
	bitsPerSample = wav.GetBitsPerSample()
	r = wav
	return
}

// Writer wav file
func (w *WAV) Writer(fileName string, channels uint16, rate uint32) (wc io.WriteCloser, err error) {
	if !strings.HasSuffix(fileName, ".wav") {
		fileName += ".wav"
	}
	file, err := os.Create(fileName)
	if err != nil {
		return
	}

	wc = wav.NewWriter(file, channels, rate, wav.S16LE)
	return
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
