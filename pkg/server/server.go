package server

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
)

type storage interface {
	List() io.ReadWriteCloser
}

type audio interface {
	Read(data []byte) (reader io.Reader, channels uint16, rate uint32, err error)
	Recode(ctx context.Context, name string, channels uint16, rate uint32, r io.Reader) (err error)
}

type udp interface {
	Send(context.Context, string, io.Reader) (err error)
	Receive(context.Context, string, io.Writer) error
}

type player interface {
	StartPlay(ctx context.Context, playerIP, playerPort, deviceName string, channels, rate uint32) (err error)
	StopPlay(ctx context.Context, playerIP, playerPort string) (err error)
}

type recoder interface {
	StartRecode(ctx context.Context, serverAddr, recoderIP, deviceName string, channels, rate int) (err error)
	StopRecode(ctx context.Context, serverAddr, recoderIP string) (err error)
}

// Server audio server
type Server struct {
	mutexPlayer  sync.Mutex
	playerClient map[string]context.CancelFunc

	mutexRecoder  sync.Mutex
	recoderClient map[string]context.CancelFunc

	audio   audio
	player  player
	recoder recoder
	storage storage
	udp     udp

	hostLayout string
}

// AddFilePlayer add player client and sening audio on client from file
func (s *Server) AddFilePlayer(ctx context.Context, playerIP, playerPort, deviceName, fileName string) (err error) {
	s.mutexPlayer.Lock()
	defer s.mutexPlayer.Unlock()

	host := fmt.Sprintf(s.hostLayout, playerIP, playerPort)
	if _, isExist := s.playerClient[host]; isExist {
		err = fmt.Errorf("%s is busy", host)
		return
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	audio, channels, rate, err := s.audio.Read(data)
	if err != nil {
		return
	}

	c, cancel := context.WithCancel(context.Background())
	if err = s.player.StartPlay(ctx, playerIP, playerPort, deviceName, uint32(channels), rate); err != nil {
		cancel()
		return
	}
	if err = s.udp.Send(c, host, audio); err != nil {
		cancel()
		s.player.StopPlay(c, playerIP, playerPort)
		return
	}

	s.playerClient[host] = cancel
	return
}

// AddFileRecoder add recoder client and receive audio signal and write on file
func (s *Server) AddFileRecoder(ctx context.Context, receivePort, fileName, recoderIP, deviceName string, channels, rate int) (err error) {
	s.mutexRecoder.Lock()
	defer s.mutexRecoder.Unlock()

	if _, isExist := s.playerClient[recoderIP]; isExist {
		err = fmt.Errorf("%s is busy", recoderIP)
		return
	}

	list := s.storage.List()
	c, cancel := context.WithCancel(context.Background())

	if err = s.udp.Receive(c, receivePort, list); err != nil {
		cancel()
		return
	}

	if err = s.audio.Recode(c, fileName, uint16(channels), uint32(rate), list); err != nil {
		cancel()
		return
	}

	if err = s.recoder.StartRecode(ctx, "TODO", recoderIP, deviceName, channels, rate); err != nil {
		cancel()
		return
	}
	s.playerClient[recoderIP] = cancel
	return
}

// AddRecoderPlayer add recoder client and player client
func (s *Server) AddRecoderPlayer(ctx context.Context, playerIP, playerPort, playerDeviceName, recoderIP, recoderDeviceName string, channels, rate int) (err error) {
	// todo
	return
}

// DeleteRecoder delete recoder client
func (s *Server) DeleteRecoder(ctx context.Context, recoderIP, receivePort string) error {
	s.mutexRecoder.Lock()
	defer s.mutexRecoder.Unlock()

	cancel, isExist := s.recoderClient[recoderIP]
	if !isExist {
		return fmt.Errorf("recoder %s not exist", recoderIP)
	}

	cancel()
	s.recoder.StopRecode(ctx, "TODO", recoderIP)
	delete(s.recoderClient, recoderIP)
	return nil
}

// DeletePlayer delete player client
func (s *Server) DeletePlayer(ctx context.Context, ip, port string) error {
	s.mutexPlayer.Lock()
	defer s.mutexPlayer.Unlock()

	host := fmt.Sprintf(s.hostLayout, ip, port)
	cancel, isExist := s.playerClient[host]
	if !isExist {
		return fmt.Errorf("player %s not exist", host)
	}

	cancel()
	s.player.StopPlay(ctx, ip, port)
	delete(s.playerClient, host)
	return nil
}

// NewServer ...
func NewServer(
	audio audio,
	recoder recoder,
	player player,
	udp udp,

	hostLayout string,
) *Server {
	return &Server{
		playerClient: make(map[string]context.CancelFunc),

		audio:   audio,
		recoder: recoder,
		player:  player,
		udp:     udp,

		hostLayout: hostLayout,
	}
}
