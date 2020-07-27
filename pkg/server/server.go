package server

import (
	"context"
)

type audio interface {
	GetSample() ([]int16, error)
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
	var (
		samples []int16
		data    []byte
	)
	if samples, err = audio.GetSample(); err != nil {
		return
	}
	for i := 0; i < len(samples)-s.size; i += s.size {
		data = s.converter.ToByte(samples[i : i+s.size])
		s.connection.Send(data)
	}
	return
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
