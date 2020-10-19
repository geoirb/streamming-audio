package playback

import (
	"context"
	"io"
	"sync"

	alsa "github.com/cocoonlife/goalsa"
)

type converter interface {
	ToInt16([]byte) []int16
}

// Playback device
type Playback struct {
	converter converter
	buffSize  int
	mutex     sync.Mutex
}

// Play audio on deviceName
func (d *Playback) Play(ctx context.Context, deviceName string, channels, rate int, r io.Reader) (err error) {
	out, err := alsa.NewPlaybackDevice(
		deviceName,
		channels,
		alsa.FormatS16LE,
		rate,
		alsa.BufferParams{},
	)
	if err != nil {
		return
	}

	go func() {
		samples := make([]byte, d.buffSize)
		for {
			select {
			case <-ctx.Done():
				out.Close()
				return
			default:
				if l, err := r.Read(samples); err == nil {
					out.Write(d.converter.ToInt16(samples[:l]))
				}
			}
		}
	}()
	return
}

// NewPlayback ...
func NewPlayback(
	converter converter,
	buffSize int,
) *Playback {
	return &Playback{
		converter: converter,
		buffSize:  buffSize,
	}
}
