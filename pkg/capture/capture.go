package capture

import (
	"context"
	"io"

	alsa "github.com/cocoonlife/goalsa"
)

type converter interface {
	ToByte([]int16) []byte
}

// Capture device
type Capture struct {
	converter converter

	buffSize int
}

// Record audio signals
func (c *Capture) Record(ctx context.Context, deviceName string, channels, rate int, dest io.WriteCloser) (err error) {
	in, err := alsa.NewCaptureDevice(
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
		defer func() {
			in.Close()
			dest.Close()
		}()
		samples := make([]int16, c.buffSize)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if n, err := in.Read(samples); err == nil {
					if _, err := dest.Write(c.converter.ToByte(samples[:n])); err != nil {
						return
					}
				}
			}
		}
	}()
	return
}

// NewCapture ..
func NewCapture(converter converter, buffSize int) *Capture {
	return &Capture{
		converter: converter,
		buffSize:  buffSize,
	}
}
