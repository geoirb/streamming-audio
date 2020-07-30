package device

import (
	"sync"
	"time"

	alsa "github.com/cocoonlife/goalsa"
)

// Device playback device
type Device struct {
	out   *alsa.PlaybackDevice
	mutex sync.Mutex

	deviceName string
	channels   int
	rate       int

	interval time.Duration
}

// Connect to device
func (d *Device) Connect() (err error) {
	d.out, err = alsa.NewPlaybackDevice(
		d.deviceName,
		d.channels,
		alsa.FormatS16LE,
		d.rate,
		alsa.BufferParams{},
	)

	d.interval = time.Duration(1 / float64(d.rate) * 1e9) 
	return
}

// Play audio track
func (d *Device) Play(audio []int16) {
	d.out.Write(audio)
	time.Sleep(d.interval)
	return
}

// Disconnect from device
func (d *Device) Disconnect() {
	d.out.Close()
}

// NewDevice ...
func NewDevice(
	deviceName string,
	channels int,
	rate int,
) *Device {
	return &Device{
		deviceName: deviceName,
		channels:   channels,
		rate:       rate,
	}
}
