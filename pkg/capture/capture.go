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
func (c *Capture) Recode(ctx context.Context, deviceName string, channels, rate int, w io.ReadWriteCloser) (err error) {
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
		samples := make([]int16, c.buffSize)
		for {
			select {
			case <-ctx.Done():
				in.Close()
				w.Close()
				return
			default:
				if n, err := in.Read(samples); err != nil && n != 0 {
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
