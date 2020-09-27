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
	StartReceive(ctx context.Context, playerIP, receivePort string) error
	StopReceive(ctx context.Context, playerIP, receivePort string) error
	StartPlay(ctx context.Context, playerIP, playerPort, deviceName string, channels, rate uint32) (err error)
	StopPlay(ctx context.Context, playerIP, deviceName string) (err error)
}

type recoder interface {
	StartRecode(ctx context.Context, serverAddr, recoderIP, deviceName string, channels, rate int) error
	StopRecode(ctx context.Context, serverAddr, recoderIP string) error
}

// Server audio server
type Server struct {
	mutexSending sync.Mutex
	sending      map[string]context.CancelFunc

	mutexPlaying sync.Mutex
	playing      map[string]context.CancelFunc

	mutexRecoder sync.Mutex
	recoding     map[string]context.CancelFunc

	audio   audio
	player  player
	recoder recoder
	storage storage
	udp     udp

	hostLayout string
	playLayout string
}

// StartSendingFile on player
func (s *Server) StartSendingFile(c context.Context, playerIP, playerPort, fileName string) (channels uint16, rate uint32, err error) {
	s.mutexSending.Lock()
	defer s.mutexSending.Unlock()

	host := fmt.Sprintf(s.hostLayout, playerIP, playerPort)
	if _, isExist := s.sending[host]; isExist {
		err = fmt.Errorf("StartSendingFile: %s is busy", host)
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

	ctx, cancel := context.WithCancel(c)
	if err = s.player.StartReceive(ctx, playerIP, playerPort); err != nil {
		cancel()
		return
	}
	if err = s.udp.Send(c, host, audio); err != nil {
		cancel()
		s.player.StopPlay(c, playerIP, playerPort)
		return
	}

	s.sending[host] = cancel
	return
}

// StopSending on player
func (s *Server) StopSending(c context.Context, playerIP, playerPort string) error {
	s.mutexSending.Lock()
	defer s.mutexSending.Unlock()

	host := fmt.Sprintf(s.hostLayout, playerIP, playerPort)
	if stop, isExist := s.sending[host]; isExist {
		stop()
		s.player.StopReceive(c, playerIP, playerPort)
		delete(s.sending, host)
		return nil
	}
	return fmt.Errorf("StopSending: %s is busy", host)
}

// StartPlaying on player
func (s *Server) StartPlaying(c context.Context, playerIP, playerPort, deviceName string, channels uint16, rate uint32) (err error) {
	s.mutexPlaying.Lock()
	defer s.mutexPlaying.Unlock()

	player := fmt.Sprintf(s.playLayout, playerIP, deviceName)
	if _, isExist := s.playing[player]; !isExist {
		ctx, cancel := context.WithCancel(c)
		if err = s.player.StartPlay(ctx, playerIP, playerPort, deviceName, uint32(channels), rate); err != nil {
			cancel()
			return
		}
		s.playing[player] = cancel
		return
	}
	err = fmt.Errorf("StartPlaying: %s is busy", player)
	return
}

// StopPlaying on player
func (s *Server) StopPlaying(c context.Context, playerIP, deviceName string) (err error) {
	s.mutexPlaying.Lock()
	defer s.mutexPlaying.Unlock()

	player := fmt.Sprintf(s.playLayout, playerIP, deviceName)
	if stop, isExist := s.playing[player]; isExist {
		stop()
		s.player.StopPlay(c, playerIP, deviceName)
		delete(s.playing, player)
		return
	}
	err = fmt.Errorf("StopPlaying: %s is not exist", player)
	return
}

// // AddFileRecoder add recoder client and receive audio signal and write on file
// func (s *Server) AddFileRecoder(ctx context.Context, receivePort, fileName, recoderIP, deviceName string, channels, rate int) (err error) {
// 	s.mutexRecoder.Lock()
// 	defer s.mutexRecoder.Unlock()

// 	if _, isExist := s.playerClient[recoderIP]; isExist {
// 		err = fmt.Errorf("%s is busy", recoderIP)
// 		return
// 	}

// 	list := s.storage.List()
// 	c, cancel := context.WithCancel(context.Background())

// 	if err = s.udp.Receive(c, receivePort, list); err != nil {
// 		cancel()
// 		return
// 	}

// 	if err = s.audio.Recode(c, fileName, uint16(channels), uint32(rate), list); err != nil {
// 		cancel()
// 		return
// 	}

// 	if err = s.recoder.StartRecode(ctx, "TODO", recoderIP, deviceName, channels, rate); err != nil {
// 		cancel()
// 		return
// 	}
// 	s.playerClient[recoderIP] = cancel
// 	return
// }

// // AddRecoderPlayer add recoder client and player client
// func (s *Server) AddRecoderPlayer(ctx context.Context, playerIP, playerPort, playerDeviceName, recoderIP, recoderDeviceName string, channels, rate int) (err error) {
// 	// todo
// 	return
// }

// // DeleteRecoder delete recoder client
// func (s *Server) DeleteRecoder(ctx context.Context, recoderIP, receivePort string) error {
// 	s.mutexRecoder.Lock()
// 	defer s.mutexRecoder.Unlock()

// 	cancel, isExist := s.recoderClient[recoderIP]
// 	if !isExist {
// 		return fmt.Errorf("recoder %s not exist", recoderIP)
// 	}

// 	cancel()
// 	s.recoder.StopRecode(ctx, "TODO", recoderIP)
// 	delete(s.recoderClient, recoderIP)
// 	return nil
// }

// NewServer ...
func NewServer(
	audio audio,
	recoder recoder,
	player player,
	udp udp,

	hostLayout string,
	playLayout string,
) *Server {
	return &Server{
		sending:  make(map[string]context.CancelFunc),
		playing:  make(map[string]context.CancelFunc),
		recoding: make(map[string]context.CancelFunc),

		audio:   audio,
		recoder: recoder,
		player:  player,
		udp:     udp,

		hostLayout: hostLayout,
		playLayout: playLayout,
	}
}
