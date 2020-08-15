package capture

import (
	"context"

	alsa "github.com/cocoonlife/goalsa"
)

type converter interface {
	ToByte([]int16) []byte
	ToInt16([]byte) []int16
}

// Capture device
type Capture struct {
	in         *alsa.CaptureDevice
	deviceName string
	channels   int
	rate       int

	sample chan []byte
	err    chan error

	converter converter
}

// Device connect
func (c *Capture) Device() (err error) {
	c.in, err = alsa.NewCaptureDevice(
		c.deviceName,
		c.channels,
		alsa.FormatS16LE,
		c.rate,
		alsa.BufferParams{},
	)

	c.sample, c.err = make(chan []byte, 1), make(chan error, 1)
	return
}

// StartReadingSamples reading audio samples
func (c *Capture) StartReadingSamples(ctx context.Context) {
	// todo
	sample := make([]int16, 4)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := c.in.Read(sample); err != nil && err != alsa.ErrOverrun {
				c.err <- err
				return
			}
			c.sample <- c.converter.ToByte(sample)
		}
	}
}

// Sample return chan for audio sample
func (c *Capture) Sample() <-chan []byte {
	return c.sample
}

// Error return chan for error
func (c *Capture) Error() <-chan error {
	return c.err
}

// StopReadingSamples ...
func (c *Capture) StopReadingSamples() {
	close(c.sample)
	close(c.err)
	c.in.Close()
}

// NewCapture ..
func NewCapture(
	deviceName string,
	channels int,
	rate int,
	converter converter,
) *Capture {
	return &Capture{
		deviceName: deviceName,
		channels:   channels,
		rate:       rate,
		converter:  converter,
	}
}
