package server

import (
	"context"
)

type audio interface {
	GetSample() (samples []int16, err error)
}

type udpConnection interface {
	Send(data []byte) (err error)
}

// Server audio server
type Server struct {
	udpConnection udpConnection
}

func (s *Server) Start(ctx context.Context, audio audio) {
	samples, _ := audio.GetSample()
}
