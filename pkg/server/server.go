package server

import (
	"context"
	"io"

	alsa "github.com/cocoonlife/goalsa"
)

type audio interface {
	GetSample() (<-chan []int16, <-chan error)
}

type connection interface {
	Send(data []byte) error
}

type converter interface {
	ToByte(src []int16) []byte
}

// Server audio server
type Server struct {
	connection connection
	converter  converter

	size int
}

// Streaming audio over connection
func (s *Server) Streaming(ctx context.Context, audio audio) (err error) {
	out, _ := alsa.NewPlaybackDevice(
		"default",
		1,
		alsa.FormatS16LE,
		44100,
		alsa.BufferParams{},
	)
	sChan, eChan := audio.GetSample()
	for {
		select {
		case samples := <-sChan:
			for i := 0; i < len(samples)-s.size; i += s.size {
				a := make([]int16, s.size)
				copy(a, samples[i:i+s.size])
				out.Write(a)
				data := s.converter.ToByte(samples[i : i+s.size])
				s.connection.Send(data)
			}
		case err = <-eChan:
			if err == io.EOF {
				return nil
			}
			return
		case <-ctx.Done():
			return
		}
	}
}

// NewServer ...
func NewServer(
	connection connection,
	converter converter,
	size int,
) *Server {
	return &Server{
		connection: connection,
		converter:  converter,
		size:       size,
	}
}
