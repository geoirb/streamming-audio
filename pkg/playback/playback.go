package playback

import (
	"sync"

	alsa "github.com/cocoonlife/goalsa"
)

// Playback device
type Playback struct {
	out        *alsa.PlaybackDevice
	deviceName string
	channels   int
	rate       int

	mutex sync.Mutex
}

// Device connect
func (d *Playback) Device() (err error) {
	d.out, err = alsa.NewPlaybackDevice(
		d.deviceName,
		d.channels,
		alsa.FormatS16LE,
		d.rate,
		alsa.BufferParams{},
	)
	return
}

// Write audio track
func (d *Playback) Write(samples []int16) {
	d.out.Write(samples)
	return
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
) *Playback {
	return &Playback{
		deviceName: deviceName,
		channels:   channels,
		rate:       rate,
	}
}
