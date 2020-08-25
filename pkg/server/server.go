package server

import (
	"context"
	"fmt"
)

type audio interface {
	StartReadingSamples(ctx context.Context, packageSize int)
	Sample() <-chan []byte
	Error() <-chan error
}

type connection interface {
	Send(data []byte) error
}

// Server audio server
type Server struct {
	pull map[connection]audio
}

// AddStreaming audio over connection
func (s *Server) AddStreaming(connection connection, audio audio) (err error) {
	if _, isExist := s.pull[connection]; isExist {
		return fmt.Errorf("connection is exist: %v", connection)
	}
	s.pull[connection] = audio
	return
}

// Start server
func (s *Server) Start(ctx context.Context) {
	for connection, audio := range s.pull {

		go s.streaming(ctx, connection, audio)
	}
}

func (s *Server) streaming(ctx context.Context, connection connection, audio audio) {
	go audio.StartReadingSamples(ctx, 1024)
	for {
		select {
		case sample := <-audio.Sample():
			if err := connection.Send(sample); err != nil {
				return
			}
		case <-audio.Error():
			return
		case <-ctx.Done():
			return
		}
	}
}

// NewServer ...
func NewServer() *Server {
	return &Server{
		pull: make(map[connection]audio),
	}
}
