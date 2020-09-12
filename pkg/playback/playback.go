package playback

import (
	"context"

	alsa "github.com/cocoonlife/goalsa"

	"github.com/geoirb/sound-ethernet-streaming/pkg/storage"
)

type converter interface {
	ToInt16([]byte) []int16
}

// Playback device
type Playback struct {
	converter converter
}

// Play audio on deviceName
func (d *Playback) Play(ctx context.Context, deviceName string, channels, rate int, s storage.List) error {
	out, err := alsa.NewPlaybackDevice(
		deviceName,
		channels,
		alsa.FormatS16LE,
		rate,
		alsa.BufferParams{},
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				out.Close()
				return
			default:
				if samples := s.Pop(); samples != nil && len(samples) > 0 {
					out.Write(d.converter.ToInt16(samples))
				}
			}
		}
	}()
	return nil
}

// NewPlayback ...
func NewPlayback(
	converter converter,
) *Playback {
	return &Playback{
		converter: converter,
	}
}
