package playback

import (
	alsa "github.com/cocoonlife/goalsa"
)

type converter interface {
	ToInt16([]byte) []int16
}

// Playback device
type Playback struct {
	out       *alsa.PlaybackDevice
	converter converter
}

// Write audio track
func (d *Playback) Write(samples []byte) {
	d.out.Write(
		d.converter.ToInt16(samples),
	)
}

// Disconnect from device
func (d *Playback) Disconnect() {
	d.out.Close()
}

// NewPlayback ...
func NewPlayback(
	deviceName string,
	channels int,
	rate int,
	converter converter,
) (p *Playback, err error) {
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

	p = &Playback{
		out:       out,
		converter: converter,
	}
	return
}
