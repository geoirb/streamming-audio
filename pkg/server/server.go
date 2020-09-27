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
	StartReceive(ctx context.Context, destIP, destPort string) (storageUUID string, err error)
	StopReceive(ctx context.Context, destIP, destPort string) (err error)
	StartPlay(ctx context.Context, playerIP, deviceName, storageUUID string, channels, rate uint32) (err error)
	StopPlay(ctx context.Context, playerIP, deviceName string) (err error)
}

type recoder interface {
	StartRecode(ctx context.Context, serverAddr, recoderIP, deviceName string, channels, rate int) error
	StopRecode(ctx context.Context, serverAddr, recoderIP string) error
}

// Server ...
type Server interface {
	StartSendingFile(ctx context.Context, destIP, destPort, fileName string) (storageUUID string, channels uint16, rate uint32, err error)
	StopSending(ctx context.Context, destIP, destPort string) (err error)
	StartPlaying(ctx context.Context, playerIP, deviceName, storageUUID string, channels uint16, rate uint32) (err error)
	StopPlaying(ctx context.Context, playerIP, deviceName string) (err error)
}

type server struct {
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
func (s *server) StartSendingFile(c context.Context, destIP, destPort, fileName string) (storageUUID string, channels uint16, rate uint32, err error) {
	s.mutexSending.Lock()
	defer s.mutexSending.Unlock()

	host := fmt.Sprintf(s.hostLayout, destIP, destPort)
	if _, isExist := s.sending[host]; isExist {
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

	ctx, cancel := context.WithCancel(c)
	if storageUUID, err = s.player.StartReceive(ctx, destIP, destPort); err != nil {
		cancel()
		return
	}
	if err = s.udp.Send(c, host, audio); err != nil {
		cancel()
		s.player.StopPlay(c, destIP, destPort)
		return
	}

	s.sending[host] = cancel
	return
}

// StopSending on player
func (s *server) StopSending(c context.Context, destIP, destPort string) error {
	s.mutexSending.Lock()
	defer s.mutexSending.Unlock()

	host := fmt.Sprintf(s.hostLayout, destIP, destPort)
	if stop, isExist := s.sending[host]; isExist {
		stop()
		s.player.StopReceive(c, destIP, destPort)
		delete(s.sending, host)
		return nil
	}
	return fmt.Errorf("%s is not exist", host)
}

// StartPlaying on player
func (s *server) StartPlaying(c context.Context, playerIP, deviceName, storageUUID string, channels uint16, rate uint32) (err error) {
	s.mutexPlaying.Lock()
	defer s.mutexPlaying.Unlock()

	player := fmt.Sprintf(s.playLayout, playerIP, deviceName)
	if _, isExist := s.playing[player]; !isExist {
		ctx, cancel := context.WithCancel(c)
		if err = s.player.StartPlay(ctx, playerIP, deviceName, storageUUID, uint32(channels), rate); err != nil {
			cancel()
			return
		}
		s.playing[player] = cancel
		return
	}
	err = fmt.Errorf("%s is busy", player)
	return
}

// StopPlaying on player
func (s *server) StopPlaying(c context.Context, playerIP, deviceName string) (err error) {
	s.mutexPlaying.Lock()
	defer s.mutexPlaying.Unlock()

	player := fmt.Sprintf(s.playLayout, playerIP, deviceName)
	if stop, isExist := s.playing[player]; isExist {
		stop()
		s.player.StopPlay(c, playerIP, deviceName)
		delete(s.playing, player)
		return
	}
	err = fmt.Errorf("%s is not exist", player)
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
) Server {
	return &server{
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
