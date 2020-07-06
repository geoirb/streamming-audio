package server

import (
	"context"
)

type source interface {
	Connect(ctx context.Context) (err error)
	Get() (src <-chan []byte)
}

type destinision interface {
	Connect(ctx context.Context) (err error)
	Send(data []byte) (err error)
}

// Server sound interface
type Server interface {
	Stream(ctx context.Context) (err error)
}

type server struct {
	src source
	dst destinision
}

func (s *server) Stream(ctx context.Context) (err error) {
	if err = s.src.Connect(ctx); err != nil {
		return
	}
	if err = s.dst.Connect(ctx); err != nil {
		return
	}

	for {
		select {
		case data := <-s.src.Get():
			if err = s.dst.Send(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// NewServer ...
func NewServer(
	src source,
	dst destinision,
) Server {
	return &server{
		src: src,
		dst: dst,
	}
}
