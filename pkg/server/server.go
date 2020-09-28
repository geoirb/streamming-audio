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
	Record(ctx context.Context, name string, channels uint16, rate uint32, r io.ReadCloser) error
}

type udp interface {
	Send(context.Context, string, io.Reader) error
	TurnOnReceiver(receivePort string) (connection io.ReadCloser, err error)
}

type player interface {
	StartReceive(ctx context.Context, destIP, destPort string) (storageUUID string, err error)
	StopReceive(ctx context.Context, destIP, destPort string) (err error)
	StartPlay(ctx context.Context, playerIP, deviceName, storageUUID string, channels, rate uint32) (err error)
	StopPlay(ctx context.Context, playerIP, deviceName string) (err error)
}

type recorder interface {
	StartRecord(ctx context.Context, destAddr, recorderIP, deviceName string, channels, rate uint32) (err error)
	StopRecord(ctx context.Context, recorderIP, deviceName string) (err error)
}

// Server ...
type Server interface {
	StartSendingFile(ctx context.Context, destIP, destPort, fileName string) (storageUUID string, channels uint16, rate uint32, err error)
	StopSending(ctx context.Context, destIP, destPort string) (err error)
	StartPlaying(ctx context.Context, playerIP, deviceName, storageUUID string, channels uint16, rate uint32) (err error)
	StopPlaying(ctx context.Context, playerIP, deviceName string) (err error)

	StartRecordingInFile(ctx context.Context, fileName, receivePort, recoderIP, deviceName string, channels, rate int) (err error)
	//todo
	StartRecordingOnPlayer(ctx context.Context, playerIP, recoderIP, deviceName string, channels, rate int) (err error)
	StopRecoding(ctx context.Context, recoderIP, deviceName string) (err error)
}

type server struct {
	mutexSending sync.Mutex
	sending      map[string]context.CancelFunc

	mutexPlaying sync.Mutex
	playing      map[string]context.CancelFunc

	mutexRecording sync.Mutex
	recoding       map[string]context.CancelFunc

	audio    audio
	player   player
	recorder recorder
	storage  storage
	udp      udp

	hostLayout   string
	deviceLayout string
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
func (s *server) StopSending(c context.Context, destIP, destPort string) (err error) {
	s.mutexSending.Lock()
	defer s.mutexSending.Unlock()

	host := fmt.Sprintf(s.hostLayout, destIP, destPort)
	if stop, isExist := s.sending[host]; isExist {
		stop()
		s.player.StopReceive(c, destIP, destPort)
		delete(s.sending, host)
		return
	}
	err = fmt.Errorf("%s is not exist", host)
	return
}

// StartPlaying on player
func (s *server) StartPlaying(c context.Context, playerIP, deviceName, storageUUID string, channels uint16, rate uint32) (err error) {
	s.mutexPlaying.Lock()
	defer s.mutexPlaying.Unlock()

	player := fmt.Sprintf(s.deviceLayout, playerIP, deviceName)
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

	player := fmt.Sprintf(s.deviceLayout, playerIP, deviceName)
	if stop, isExist := s.playing[player]; isExist {
		stop()
		s.player.StopPlay(c, playerIP, deviceName)
		delete(s.playing, player)
		return
	}
	err = fmt.Errorf("%s is not exist", player)
	return
}

// StartRecordingInFile start recoding player and save on file
func (s *server) StartRecordingInFile(c context.Context, fileName, receivePort, recoderIP, deviceName string, channels, rate int) (err error) {
	s.mutexRecording.Lock()
	defer s.mutexRecording.Unlock()

	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, deviceName)
	if _, isExist := s.recoding[recoder]; !isExist {
		var connection io.ReadCloser
		if connection, err = s.udp.TurnOnReceiver(receivePort); err == nil {
			ctx, stop := context.WithCancel(context.Background())
			if err = s.audio.Record(ctx, fileName, uint16(channels), uint32(rate), connection); err == nil {
				//todo
				if err = s.recorder.StartRecord(ctx, "127.0.0.1:"+receivePort, recoderIP, deviceName, uint32(channels), uint32(rate)); err == nil {
					s.recoding[recoder] = stop
					return
				}
			}
			stop()
		}
		return
	}
	err = fmt.Errorf("%s is busy", recoder)
	return
}

// StartRecordingInFile start recoding player and save on file
// todo
func (s *server) StartRecordingOnPlayer(c context.Context, playerIP, recoderIP, deviceName string, channels, rate int) (storageUUID string, err error) {
	s.mutexRecording.Lock()
	defer s.mutexRecording.Unlock()

	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, deviceName)
	if _, isExist := s.recoding[recoder]; !isExist {
		ctx, stop := context.WithCancel(context.Background())
		if storageUUID, err = s.player.StartReceive(ctx, destIP, destPort); err != nil {
			stop()
			return
		}
		if err = s.recorder.StartRecord(ctx, playerIP, recoderIP, deviceName, uint32(channels), uint32(rate)); err == nil {
			s.recoding[recoder] = stop
			return
		}
		stop()
		return
	}
	err = fmt.Errorf("%s is busy", recoder)
	return
}

func (s *server) StopRecoding(c context.Context, recoderIP, deviceName string) (err error) {
	s.mutexRecording.Lock()
	defer s.mutexRecording.Unlock()

	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, deviceName)
	if stop, isExist := s.recoding[recoder]; isExist {
		stop()
		s.recorder.StopRecord(c, recoderIP, deviceName)
		delete(s.recoding, recoder)
		return
	}
	err = fmt.Errorf("%s is not exist", recoder)
	return
}

// NewServer ...
func NewServer(
	audio audio,
	recorder recorder,
	player player,
	udp udp,

	hostLayout string,
	deviceLayout string,
) Server {
	return &server{
		sending:  make(map[string]context.CancelFunc),
		playing:  make(map[string]context.CancelFunc),
		recoding: make(map[string]context.CancelFunc),

		audio:    audio,
		recorder: recorder,
		player:   player,
		udp:      udp,

		hostLayout:   hostLayout,
		deviceLayout: deviceLayout,
	}
}
