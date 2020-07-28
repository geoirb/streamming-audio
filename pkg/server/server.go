package server

import (
	"context"
	"io"
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
	sChan, eChan := audio.GetSample()
	var data []byte
	for {
		select {
		case samples := <-sChan:
			for i := 0; i < len(samples)-s.size; i += s.size {
				data = s.converter.ToByte(samples[i : i+s.size])
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
