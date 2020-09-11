package server

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
)

type file interface {
	Audio(data []byte) (reader io.Reader, channels uint16, rate uint32, err error)
}

type udp interface {
	Send(context.Context, string, io.Reader) (err error)
}

type media interface {
	StartReceive(ctx context.Context, ip, port, deviceName string, channels, rate uint32) (err error)
	StopReceive(ctx context.Context, ip, port string) (err error)
}

// Server audio server
type Server struct {
	hostLayout string // "%s:%d"

	file  file
	media media
	udp   udp

	mutex  sync.Mutex
	client map[string]context.CancelFunc
}

// AddFileMedia add media client and sening audio on client from file
func (s *Server) AddFileMedia(ctx context.Context, ip, port, deviceName, fileName string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	host := fmt.Sprintf(s.hostLayout, ip, port)
	if _, isExist := s.client[host]; isExist {
		err = fmt.Errorf("%s is busy", host)
		return
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	audio, channels, rate, err := s.file.Audio(data)
	if err != nil {
		return
	}

	c, cancel := context.WithCancel(ctx)
	if err = s.media.StartReceive(c, ip, port, deviceName, uint32(channels), rate); err != nil {
		cancel()
		return
	}
	if err = s.udp.Send(c, host, audio); err != nil {
		cancel()
		s.media.StopReceive(c, ip, port)
		return
	}

	s.client[host] = cancel
	return
}

// DeleteMedia delete media client
func (s *Server) DeleteMedia(ctx context.Context, ip, port string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	host := fmt.Sprintf(s.hostLayout, ip, port)
	cancel, isExist := s.client[host]
	if !isExist {
		err = fmt.Errorf("client %s not exist", host)
		return
	}

	cancel()
	s.media.StopReceive(ctx, ip, port)
	return
}

// NewServer ...
func NewServer(
	hostLayout string,
	file file,
	media media,
	udp udp,
) *Server {
	return &Server{
		hostLayout: hostLayout,
		file:       file,
		media:      media,
		udp:        udp,
		client:     make(map[string]context.CancelFunc),
	}
}
