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

type player interface {
	StartPlay(ctx context.Context, ip, port, deviceName string, channels, rate uint32) (err error)
	StopPlay(ctx context.Context, ip, port string) (err error)
}

// Server audio server
type Server struct {
	mutex  sync.Mutex
	client map[string]context.CancelFunc

	file   file
	player player
	udp    udp

	hostLayout string
}

// AddFilePlayer add player client and sening audio on client from file
func (s *Server) AddFilePlayer(ctx context.Context, ip, port, deviceName, fileName string) (err error) {
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
	if err = s.player.StartPlay(c, ip, port, deviceName, uint32(channels), rate); err != nil {
		cancel()
		return
	}
	if err = s.udp.Send(c, host, audio); err != nil {
		cancel()
		s.player.StopPlay(c, ip, port)
		return
	}

	s.client[host] = cancel
	return
}

// DeletePlayer delete player client
func (s *Server) DeletePlayer(ctx context.Context, ip, port string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	host := fmt.Sprintf(s.hostLayout, ip, port)
	cancel, isExist := s.client[host]
	if !isExist {
		err = fmt.Errorf("client %s not exist", host)
		return
	}

	cancel()
	s.player.StopPlay(ctx, ip, port)
	return
}

// NewServer ...
func NewServer(
	file file,
	player player,
	udp udp,

	hostLayout string,
) *Server {
	return &Server{
		client: make(map[string]context.CancelFunc),

		file:   file,
		player: player,
		udp:    udp,

		hostLayout: hostLayout,
	}
}
