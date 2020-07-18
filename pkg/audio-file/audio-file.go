package audiofile

import (
	"context"
	"io/ioutil"
)

// AudioFile interface for read with audio file
type AudioFile interface {
	Connect(ctx context.Context) (src <-chan []byte, err error)
}

type audioFile struct {
	fileName string
	src      chan []byte
}

func (a *audioFile) Connect(ctx context.Context) (src <-chan []byte, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(a.fileName); err != nil {
		return
	}

	go func(data []byte) {
		a.src <- data
	}(data)

	src = a.src
	return
}

// NewAudioFile ...
func NewAudioFile(
	fileName string,
) AudioFile {
	return &audioFile{
		fileName: fileName,

		src: make(chan []byte, 1),
	}
}
