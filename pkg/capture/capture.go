package capture

import (
	alsa "github.com/cocoonlife/goalsa"
)

type converter interface {
	ToByte([]int16) []byte
	ToInt16([]byte) []int16
}

// Capture device
type Capture struct {
	in *alsa.CaptureDevice

	converter converter
}

// Read audio bytes
func (c *Capture) Read() ([]byte, error) {
	//todo
	samples := make([]int16, 8)
	l, err := c.in.Read(samples)
	return c.converter.ToByte(samples[:l]), err
}

// Disconnect from capture device
func (c *Capture) Disconnect() {
	c.in.Close()
}

// NewCapture ..
func NewCapture(
	deviceName string,
	channels int,
	rate int,
	converter converter,
) (c *Capture, err error) {
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

	c = &Capture{
		in:        in,
		converter: converter,
	}
	return
}
