package wav

import (
	"bytes"
	"fmt"
	"io"

	"github.com/cryptix/wav"
)

// WAV audio file
type WAV struct {
	reader *wav.Reader
}

// Parse wav file
func (w *WAV) Parse(data []byte) (err error) {
	f := bytes.NewReader(data)
	w.reader, err = wav.NewReader(
		f,
		f.Size(),
	)
	return
}

// GetSample get audio sample
func (w *WAV) GetSample(size int) (<-chan []byte, <-chan error) {
	c, e := make(chan []byte), make(chan error)
	if w.reader == nil {
		e <- fmt.Errorf("wav info not exist")
		return c, e
	}

	b := size % int(w.reader.GetBitsPerSample()/8)
	if b != 0 {
		e <- fmt.Errorf("wav info not exist")
		return c, e
	}

	reader, err := w.reader.GetDumbReader()
	if err != nil {
		e <- err
		return c, e
	}

	go func() {
		defer func() {
			close(c)
			close(e)
		}()

		s := make([]byte, size)
		for {
			if _, err = reader.Read(s); err == io.EOF {
				return
			}

			if err != nil {
				e <- err
				return
			}
			c <- s
		}
	}()
	return c, e
}

// NewWAV return handler wav file
func NewWAV() *WAV {
	return &WAV{}
}
