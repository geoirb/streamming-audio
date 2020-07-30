package server

import (
	"context"
)

type audio interface {
	GetSample(size int) (<-chan []byte, <-chan error)
}

type connection interface {
	Send(data []byte) error
}

type converter interface {
	ToInt16(src []byte) (dst []int16)
}

// Server audio server
type Server struct {
	connection connection
	converter  converter

	size int
}

// Streaming audio over connection
func (s *Server) Streaming(ctx context.Context, audio audio) (err error) {
	sChan, eChan := audio.GetSample(s.size)
	for {
		select {
		case samples := <-sChan:
			s.connection.Send(samples)
		case err = <-eChan:
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
