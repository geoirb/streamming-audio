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

// Recode audio signals
func (c *Capture) Recode(ctx context.Context, deviceName string, channels, rate int, w io.Writer) (err error) {
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

	samples := make([]int16, c.buffSize)
	go func() {
		for {
			select {
			case <-ctx.Done():
				in.Close()
				return
			default:
				if _, err := in.Read(samples); err != nil {
					w.Write(c.converter.ToByte(samples))
				}
			}
		}
	}()
	return
}

// NewCapture ..
func NewCapture() *Capture {
	return &Capture{}
}
