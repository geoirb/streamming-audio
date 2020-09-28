package capture

import (
	"context"
	"fmt"
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
		<-ctx.Done()
		in.Close()
		dest.Close()
	}()

	go func() {
		samples := make([]int16, c.buffSize)
		for {
			n, _ := in.Read(samples)
			fmt.Println(samples)
			if _, err := dest.Write(c.converter.ToByte(samples[:n])); err != nil {
				return
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
