package server

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
)

type audio interface {
	Read(data []byte) (reader io.Reader, channels uint16, rate uint32, err error)
	Record(ctx context.Context, name string, channels uint16, rate uint32, r io.ReadCloser) error
}

type udp interface {
	Send(context.Context, string, io.Reader) error
	TurnOnReceiver(receivePort string) (connection io.ReadCloser, err error)
}

type player interface {
	StartReceive(ctx context.Context, playerIP, receivePort string, UUID *string) (storageUUID string, err error)
	StopReceive(ctx context.Context, destIP, destPort string) (err error)
	StartPlay(ctx context.Context, playerIP, deviceName, storageUUID string, channels, rate uint32) (err error)
	StopPlay(ctx context.Context, playerIP, deviceName string) (err error)
	ClearStorage(ctx context.Context, playerIP, storageUUID string) (err error)
}

type recorder interface {
	StartRecord(ctx context.Context, destAddr, recorderIP, deviceName string, channels, rate uint32) (err error)
	StopRecord(ctx context.Context, recorderIP, deviceName string) (err error)
}

// Server ...
type Server interface {
	PlayAudioFile(ctx context.Context, playerIP, playerPort, fileName, deviceName string) (storageUUID string, channels uint16, rate uint32, err error)
	Play(ctx context.Context, playerIP, storageUUID, deviceName string, channels uint16, rate uint32) (err error)
	Pause(ctx context.Context, playerIP, deviceName string) (err error)
	Stop(c context.Context, playerIP, playerPort, deviceName, storageUUID string) (err error)

	StartRecordingOnPlayer(ctx context.Context, playerIP, playerPort, playerDeviceName, recoderIP, recorderDeviceName string, channels, rate int) (storageUUID string, err error)
}

type server struct {
	mutexSending sync.RWMutex
	sending      map[string]context.CancelFunc

	mutexPlaying sync.RWMutex
	playing      map[string]struct{}

	mutexRecording sync.RWMutex
	recoding       map[string]struct{}

	audio    audio
	player   player
	recorder recorder
	udp      udp

	hostLayout   string
	deviceLayout string
}

func (s *server) PlayAudioFile(c context.Context, playerIP, playerPort, fileName, deviceName string) (storageUUID string, channels uint16, rate uint32, err error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}

	audio, channels, rate, err := s.audio.Read(data)
	if err != nil {
		return
	}

	if storageUUID, err = s.startSending(c, playerIP, playerPort, audio); err == nil {
		err = s.Play(c, playerIP, storageUUID, deviceName, channels, rate)
	}
	return
}

// Play on player
func (s *server) Play(c context.Context, playerIP, storageUUID, deviceName string, channels uint16, rate uint32) (err error) {
	s.mutexPlaying.RLock()
	player := fmt.Sprintf(s.deviceLayout, playerIP, deviceName)
	if _, isExist := s.playing[player]; !isExist {
		s.mutexPlaying.RUnlock()

		if err = s.player.StartPlay(c, playerIP, deviceName, storageUUID, uint32(channels), rate); err == nil {
			s.mutexPlaying.Lock()
			s.playing[player] = struct{}{}
			s.mutexPlaying.Unlock()
		}
		return
	}
	s.mutexPlaying.RUnlock()
	err = fmt.Errorf("%s is busy", player)
	return
}

// Stop on player
func (s *server) Pause(c context.Context, playerIP, deviceName string) (err error) {
	s.mutexPlaying.RLock()
	player := fmt.Sprintf(s.deviceLayout, playerIP, deviceName)
	if _, isExist := s.playing[player]; isExist {
		s.mutexPlaying.RUnlock()

		if err = s.player.StopPlay(c, playerIP, deviceName); err == nil {
			s.mutexPlaying.Lock()
			delete(s.playing, player)
			s.mutexPlaying.Unlock()
		}
		return
	}
	s.mutexPlaying.RUnlock()
	err = fmt.Errorf("%s is not exist", player)
	return
}

// Stop on player
func (s *server) Stop(c context.Context, playerIP, playerPort, deviceName, storageUUID string) (err error) {
	if err = s.stopSending(c, playerIP, playerPort); err != nil {
		return
	}

	s.mutexPlaying.RLock()
	player := fmt.Sprintf(s.deviceLayout, playerIP, deviceName)
	if _, isExist := s.playing[player]; isExist {
		s.mutexPlaying.RUnlock()

		if err = s.player.StopPlay(c, playerIP, deviceName); err != nil {
			return
		}
		if err = s.player.ClearStorage(c, playerIP, storageUUID); err != nil {
			return
		}

		s.mutexPlaying.Lock()
		delete(s.playing, player)
		s.mutexPlaying.Unlock()
		return
	}
	s.mutexPlaying.RUnlock()
	err = fmt.Errorf("%s is not exist", player)
	return
}

// StartRecordingInFile start recoding player and save on file
// todo
func (s *server) StartRecordingOnPlayer(c context.Context, playerIP, playerPort, playerDeviceName, recoderIP, recorderDeviceName string, channels, rate int) (storageUUID string, err error) {
	s.mutexRecording.RLock()

	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, recorderDeviceName)
	if _, isExist := s.recoding[recoder]; !isExist {
		s.mutexRecording.RUnlock()

		if storageUUID, err = s.player.StartReceive(c, playerIP, playerPort, nil); err != nil {
			return
		}
		if err = s.recorder.StartRecord(c, playerIP+":"+playerPort, recoderIP, recorderDeviceName, uint32(channels), uint32(rate)); err != nil {
			s.player.StopReceive(c, playerIP, playerPort)
			s.player.ClearStorage(c, playerIP, storageUUID)
			return
		}
		if err = s.Play(c, playerIP, storageUUID, playerDeviceName, uint16(channels), uint32(rate)); err != nil {
			s.player.StopReceive(c, playerIP, playerPort)
			s.player.ClearStorage(c, playerIP, storageUUID)
			s.recorder.StopRecord(c, recoderIP, recorderDeviceName)
			return
		}
		s.mutexRecording.Lock()
		s.recoding[recoder] = struct{}{}
		s.mutexRecording.Unlock()
		return
	}
	s.mutexRecording.RUnlock()
	err = fmt.Errorf("%s is busy", recoder)
	return
}

func (s *server) startSending(c context.Context, destIP, destPort string, r io.Reader) (storageUUID string, err error) {
	s.mutexSending.RLock()
	host := fmt.Sprintf(s.hostLayout, destIP, destPort)
	if _, isExist := s.sending[host]; isExist {
		s.mutexSending.RUnlock()
		err = fmt.Errorf("%s is busy", host)
		return
	}
	s.mutexSending.RUnlock()

	ctx, cancel := context.WithCancel(context.Background())
	if storageUUID, err = s.player.StartReceive(c, destIP, destPort, nil); err != nil {
		cancel()
		return
	}
	if err = s.udp.Send(ctx, host, r); err == nil {
		s.mutexSending.Lock()
		s.sending[host] = cancel
		s.mutexSending.Unlock()
		return
	}

	cancel()
	s.player.StopPlay(c, destIP, destPort)
	return
}

// StopSending on player
func (s *server) stopSending(c context.Context, destIP, destPort string) (err error) {
	s.mutexSending.RLock()
	host := fmt.Sprintf(s.hostLayout, destIP, destPort)
	if stop, isExist := s.sending[host]; isExist {
		s.mutexSending.RUnlock()

		stop()
		if err = s.player.StopReceive(c, destIP, destPort); err == nil {
			s.mutexSending.Lock()
			delete(s.sending, host)
			s.mutexSending.Unlock()
		}
		return
	}
	s.mutexSending.RUnlock()
	err = fmt.Errorf("%s is not exist", host)
	return
}

// // StartRecordingInFile start recoding player and save on file
// func (s *server) StartRecordingInFile(c context.Context, fileName, receivePort, recoderIP, deviceName string, channels, rate int) (err error) {
// 	s.mutexRecording.Lock()
// 	defer s.mutexRecording.Unlock()

// 	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, deviceName)
// 	if _, isExist := s.recoding[recoder]; !isExist {
// 		var connection io.ReadCloser
// 		if connection, err = s.udp.TurnOnReceiver(receivePort); err == nil {
// 			ctx, stop := context.WithCancel(context.Background())
// 			if err = s.audio.Record(ctx, fileName, uint16(channels), uint32(rate), connection); err == nil {
// 				//todo
// 				if err = s.recorder.StartRecord(ctx, "127.0.0.1:"+receivePort, recoderIP, deviceName, uint32(channels), uint32(rate)); err == nil {
// 					s.recoding[recoder] = stop
// 					return
// 				}
// 			}
// 			stop()
// 		}
// 		return
// 	}
// 	err = fmt.Errorf("%s is busy", recoder)
// 	return
// }

func (s *server) StopRecoding(c context.Context, recoderIP, deviceName string) (err error) {
	s.mutexRecording.RLock()

	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, deviceName)
	if _, isExist := s.recoding[recoder]; isExist {
		s.mutexRecording.RUnlock()

		s.recorder.StopRecord(c, recoderIP, deviceName)

		s.mutexRecording.Lock()
		delete(s.recoding, recoder)
		s.mutexRecording.Unlock()
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
		playing:  make(map[string]struct{}),
		recoding: make(map[string]struct{}),

		audio:    audio,
		recorder: recorder,
		player:   player,
		udp:      udp,

		hostLayout:   hostLayout,
		deviceLayout: deviceLayout,
	}
}
