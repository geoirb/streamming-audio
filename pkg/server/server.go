package server

import (
	"context"
	"fmt"
)

type audio interface {
	Read() ([]byte, error)
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
	for {
		select {
		case <-ctx.Done():
			return
		default:
			samples, err := audio.Read()
			if err != nil {
				return
			}
			connection.Send(samples)
		}
	}
}

// NewServer ...
func NewServer() *Server {
	return &Server{
		pull: make(map[connection]audio),
	}
}
