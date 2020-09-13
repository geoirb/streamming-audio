package playback

import (
	"context"
	"io"

	alsa "github.com/cocoonlife/goalsa"
)

type converter interface {
	ToInt16([]byte) []int16
}

// Playback device
type Playback struct {
	converter converter

	buffSize int
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

	samples := make([]byte, d.buffSize)
	go func() {
		for {
			select {
			case <-ctx.Done():
				out.Close()
				return
			default:
				if _, err := r.Read(samples); err == nil {
					out.Write(d.converter.ToInt16(samples))
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

		buffSize: buffSize,
	}
}
