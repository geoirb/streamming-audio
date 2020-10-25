package playback

import (
	"context"
	"io"

	alsa "github.com/cocoonlife/goalsa"
)

var formatList map[int]alsa.Format = map[int]alsa.Format{
	16: alsa.FormatS16LE,
	24: alsa.FormatS24LE,
	32: alsa.FormatS32LE,
}

type converter interface {
	ToInt16([]byte) []int16
}

// Playback device
type Playback struct {
	converter converter
	buffSize  int
}

// Play audio on deviceName
func (d *Playback) Play(ctx context.Context, deviceName string, channels, rate, bitsPerSample int, r io.Reader) (err error) {
	// format, isExist := formatList[bitsPerSample]
	// if !isExist {
	// 	err = ErrFormatNotExist
	// 	return
	// }

	out, err := alsa.NewPlaybackDevice(
		deviceName,
		channels,
		alsa.FormatS24LE,
		rate,
		alsa.BufferParams{},
	)
	if err != nil {
		return
	}

	go func() {
		<-ctx.Done()
		out.Close()
	}()

	go func() {
		samples := make([]byte, d.buffSize)
		for {
			if l, err := r.Read(samples); err == nil {
				out.Write(d.converter.ToInt16(samples[:l]))
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
