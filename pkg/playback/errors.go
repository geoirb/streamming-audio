package playback

import (
	"errors"
)

var (
	// ErrFormatNotExist not exist format for alsa playback device
	ErrFormatNotExist = errors.New("format for alsa not exist")
)
