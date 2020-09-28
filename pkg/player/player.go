package player

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/twinj/uuid"
)

type storage interface {
	List() io.ReadWriteCloser
}

type udp interface {
	Receive(context.Context, string, io.Writer) error
}

type device interface {
	Play(context.Context, string, int, int, io.Reader) error
}

type playing struct {
	stop        context.CancelFunc
	storageUUID string
}

// Player audio signal
type Player struct {
	receivingMutex sync.Mutex
	receiving      map[string]context.CancelFunc

	storagingMutex sync.Mutex
	storaging      map[string]io.ReadWriteCloser

	playingMutex sync.Mutex
	playing      map[string]playing

	udp     udp
	device  device
	storage storage
}

// StartReceive start receive data from server and save
func (p *Player) StartReceive(c context.Context, in *StartReceiveRequest) (out *StartReceiveResponse, err error) {
	p.receivingMutex.Lock()
	defer p.receivingMutex.Unlock()

	if _, isExist := p.receiving[in.Port]; !isExist {
		storage := p.storage.List()
		ctx, stop := context.WithCancel(context.Background())
		if err = p.udp.Receive(ctx, in.Port, storage); err == nil {
			uuid := uuid.NewV4().String()

			p.storagingMutex.Lock()
			p.storaging[uuid] = storage
			p.storagingMutex.Unlock()
			
			p.receiving[in.Port] = stop
			
			out = &StartReceiveResponse{
				StorageUUID: uuid,
			}
			return
		}
		stop()
		return
	}
	err = fmt.Errorf("%v is busy", in.Port)
	return
}

// StopReceive stop receive data from server
// Force in StopReceiveRequest delete storage
func (p *Player) StopReceive(c context.Context, in *StopReceiveRequest) (out *StopReceiveResponse, err error) {
	p.receivingMutex.Lock()
	defer p.receivingMutex.Unlock()

	if stop, isExist := p.receiving[in.Port]; isExist {
		stop()
		delete(p.receiving, in.Port)
		out = &StopReceiveResponse{}
		return
	}
	err = fmt.Errorf("%v is not exist", in.Port)
	return
}

// StartPlay play audio on device
func (p *Player) StartPlay(c context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	p.storagingMutex.Lock()
	storage, isExist := p.storaging[in.StorageUUID]
	if !isExist {
		p.storagingMutex.Unlock()
		err = fmt.Errorf("storage %v is not exist", in.StorageUUID)
		return
	}
	p.storagingMutex.Unlock()

	p.playingMutex.Lock()
	defer p.playingMutex.Unlock()

	if _, isExist := p.playing[in.DeviceName]; !isExist {
		ctx, stop := context.WithCancel(context.Background())
		if err = p.device.Play(ctx, in.DeviceName, int(in.Channels), int(in.Rate), storage); err == nil {
			p.playing[in.DeviceName] = playing{
				stop:        stop,
				storageUUID: in.StorageUUID,
			}
			out = &StartPlayResponse{}
			return
		}
		stop()
		return
	}
	err = fmt.Errorf("%s is busy", in.DeviceName)
	return
}

// StopPlay stop play on device
func (p *Player) StopPlay(c context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
	p.playingMutex.Lock()
	defer p.playingMutex.Unlock()

	if playing, isExist := p.playing[in.DeviceName]; isExist {
		playing.stop()
		delete(p.playing, in.DeviceName)

		p.storagingMutex.Lock()
		delete(p.playing, playing.storageUUID)
		p.storagingMutex.Unlock()

		out = &StopPlayResponse{}
		return
	}
	err = fmt.Errorf("%s is not exist", in.DeviceName)
	return
}

// NewPlayer ...
func NewPlayer(
	udp udp,
	device device,
	storage storage,
) PlayerServer {
	return &Player{
		receiving: make(map[string]context.CancelFunc),
		storaging: make(map[string]io.ReadWriteCloser),
		playing:   make(map[string]playing),

		udp:     udp,
		device:  device,
		storage: storage,
	}
}
