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
	Write(name string, channels uint16, rate uint32) (io.WriteCloser, error)
}

type udp interface {
	Send(context.Context, string, io.Reader) error
	Receive(ctx context.Context, receivePort string, w io.Writer) (err error)
}

type player interface {
	StartReceive(ctx context.Context, playerIP, receivePort string, UUID *string) (storageUUID string, err error)
	StopReceive(ctx context.Context, destIP, destPort string) (err error)
	StartPlay(ctx context.Context, playerIP, deviceName, storageUUID string, channels, rate uint32) (err error)
	StopPlay(ctx context.Context, playerIP, deviceName string) (err error)
	ClearStorage(ctx context.Context, playerIP, storageUUID string) (err error)
}

type recorder interface {
	StartSend(ctx context.Context, destAddr, recorderIP, deviceName string, channels, rate uint32) (err error)
	StopSend(ctx context.Context, recorderIP, deviceName string) (err error)
}

// Server ...
type Server interface {
	PlayAudioFile(ctx context.Context, playerIP, playerPort, fileName, deviceName string) (storageUUID string, channels uint16, rate uint32, err error)
	Play(ctx context.Context, playerIP, storageUUID, deviceName string, channels uint16, rate uint32) (err error)
	Pause(ctx context.Context, playerIP, deviceName string) (err error)
	Stop(c context.Context, playerIP, playerPort, deviceName, storageUUID string) (err error)

	RecordingOnPlayer(ctx context.Context, playerIP, playerPort, playerDeviceNameName, recoderIP, recorderDeviceName string, channels, rate int) (storageUUID string, err error)
	RecordingInFile(c context.Context, fileName, receivePort, recoderIP, recoderDeviceName string, channels, rate int) (err error)
	StopRecoding(c context.Context, recoderIP, deviceName string) (err error)
}

type server struct {
	mutexSending sync.Mutex
	sending      map[string]context.CancelFunc

	mutexPlaying sync.Mutex
	playing      map[string]struct{}

	mutexRecording sync.Mutex
	recoding       map[string]context.CancelFunc

	audio    audio
	player   player
	recorder recorder
	udp      udp

	hostLayout   string
	deviceLayout string
}

// PlayAudioFile send file on player client and play then
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

// Play audio from storage on player client
func (s *server) Play(c context.Context, playerIP, storageUUID, deviceName string, channels uint16, rate uint32) (err error) {
	s.mutexPlaying.Lock()
	defer s.mutexPlaying.Unlock()

	player := fmt.Sprintf(s.deviceLayout, playerIP, deviceName)
	if _, isExist := s.playing[player]; !isExist {
		if err = s.player.StartPlay(c, playerIP, deviceName, storageUUID, uint32(channels), rate); err == nil {
			s.playing[player] = struct{}{}
		}
		return
	}
	err = fmt.Errorf("%s is busy", player)
	return
}

//  Pause audio from storage on player client
func (s *server) Pause(c context.Context, playerIP, deviceName string) (err error) {
	s.mutexPlaying.Lock()
	defer s.mutexPlaying.Unlock()

	player := fmt.Sprintf(s.deviceLayout, playerIP, deviceName)
	if _, isExist := s.playing[player]; isExist {
		if err = s.player.StopPlay(c, playerIP, deviceName); err == nil {
			delete(s.playing, player)
		}
		return
	}
	err = fmt.Errorf("%s is not exist", player)
	return
}

// Stop audio from storage on player client
func (s *server) Stop(c context.Context, playerIP, playerPort, deviceName, storageUUID string) (err error) {
	if err = s.stopSending(c, playerIP, playerPort); err != nil {
		return
	}

	s.mutexPlaying.Lock()
	defer s.mutexPlaying.Unlock()

	player := fmt.Sprintf(s.deviceLayout, playerIP, deviceName)
	if _, isExist := s.playing[player]; isExist {
		if err = s.player.StopPlay(c, playerIP, deviceName); err != nil {
			return
		}
		if err = s.player.ClearStorage(c, playerIP, storageUUID); err != nil {
			return
		}
		delete(s.playing, player)
		return
	}
	err = fmt.Errorf("%s is not exist", player)
	return
}

// StartSendingInFile start recoding player and save on file
func (s *server) RecordingOnPlayer(c context.Context, playerIP, playerPort, playerDeviceNameName, recoderIP, recorderDeviceName string, channels, rate int) (storageUUID string, err error) {
	s.mutexRecording.Lock()
	defer s.mutexRecording.Unlock()

	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, recorderDeviceName)
	if _, isExist := s.recoding[recoder]; !isExist {
		if storageUUID, err = s.player.StartReceive(c, playerIP, playerPort, nil); err != nil {
			return
		}
		if err = s.recorder.StartSend(c, playerIP+":"+playerPort, recoderIP, recorderDeviceName, uint32(channels), uint32(rate)); err != nil {
			s.player.StopReceive(c, playerIP, playerPort)
			s.player.ClearStorage(c, playerIP, storageUUID)
			return
		}
		if err = s.Play(c, playerIP, storageUUID, playerDeviceNameName, uint16(channels), uint32(rate)); err != nil {
			s.player.StopReceive(c, playerIP, playerPort)
			s.player.ClearStorage(c, playerIP, storageUUID)
			s.recorder.StopSend(c, recoderIP, recorderDeviceName)
			return
		}
		var stop context.CancelFunc
		s.recoding[recoder] = stop
		return
	}
	err = fmt.Errorf("%s is busy", recoder)
	return
}

// StartSendingInFile start recoding player and save on file
func (s *server) RecordingInFile(c context.Context, fileName, receivePort, recoderIP, recoderDeviceName string, channels, rate int) (err error) {
	s.mutexRecording.Lock()
	defer s.mutexRecording.Unlock()

	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, recoderDeviceName)
	if _, isExist := s.recoding[recoder]; !isExist {
		var writer io.WriteCloser
		ctx, stop := context.WithCancel(context.Background())
		if writer, err = s.audio.Write(fileName, uint16(channels), uint32(rate)); err == nil {
			if err = s.udp.Receive(ctx, receivePort, writer); err == nil {
				if err = s.recorder.StartSend(c, fmt.Sprintf(s.hostLayout, "127.0.0.1", receivePort), recoderIP, recoderDeviceName, uint32(channels), uint32(rate)); err == nil {
					s.recoding[recoder] = func() {
						stop()
						writer.Close()
					}
					return
				}
			}
		}
		stop()
		return
	}
	err = fmt.Errorf("%s is busy", recoder)
	return
}

// StopSending ...
func (s *server) StopRecoding(c context.Context, recoderIP, deviceName string) (err error) {
	s.mutexRecording.Lock()
	defer s.mutexRecording.Unlock()

	recoder := fmt.Sprintf(s.deviceLayout, recoderIP, deviceName)
	if _, isExist := s.recoding[recoder]; isExist {
		s.recorder.StopSend(c, recoderIP, deviceName)
		delete(s.recoding, recoder)
		return
	}
	err = fmt.Errorf("%s is not exist", recoder)
	return
}

func (s *server) startSending(c context.Context, destIP, destPort string, r io.Reader) (storageUUID string, err error) {
	s.mutexSending.Lock()
	defer s.mutexSending.Unlock()

	host := fmt.Sprintf(s.hostLayout, destIP, destPort)
	if _, isExist := s.sending[host]; isExist {
		err = fmt.Errorf("%s is busy", host)
		return
	}

	ctx, stop := context.WithCancel(context.Background())
	if storageUUID, err = s.player.StartReceive(c, destIP, destPort, nil); err != nil {
		stop()
		return
	}
	if err = s.udp.Send(ctx, host, r); err == nil {
		s.sending[host] = stop
		return
	}

	stop()
	s.player.StopPlay(c, destIP, destPort)
	return
}

// StopSending on player
func (s *server) stopSending(c context.Context, destIP, destPort string) (err error) {
	s.mutexSending.Lock()
	defer s.mutexSending.Unlock()

	host := fmt.Sprintf(s.hostLayout, destIP, destPort)
	if stop, isExist := s.sending[host]; isExist {
		stop()
		if err = s.player.StopReceive(c, destIP, destPort); err == nil {
			delete(s.sending, host)
		}
		return
	}
	err = fmt.Errorf("%s is not exist", host)
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
		recoding: make(map[string]context.CancelFunc),

		audio:    audio,
		recorder: recorder,
		player:   player,
		udp:      udp,

		hostLayout:   hostLayout,
		deviceLayout: deviceLayout,
	}
}
