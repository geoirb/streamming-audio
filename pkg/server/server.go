package server

import (
	"context"
	"fmt"
)

type audio interface {
	StartReadingSample(ctx context.Context)
	Sample() <-chan []byte
	Error() <-chan error
	StopReadingSample()
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
	go audio.StartReadingSample(ctx)
	defer audio.StopReadingSample()
	for {
		select {
		case sample := <-audio.Sample():
			connection.Send(sample)
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
