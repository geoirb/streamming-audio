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

type player struct {
	receivingMutex sync.RWMutex
	receiving      map[string]context.CancelFunc

	storagingMutex sync.RWMutex
	storaging      map[string]io.ReadWriteCloser

	playingMutex sync.RWMutex
	playing      map[string]context.CancelFunc

	udp     udp
	device  device
	storage storage
}

// StartReceive start receive data from server and save
func (p *player) StartReceive(c context.Context, in *StartReceiveRequest) (out *StartReceiveResponse, err error) {
	p.receivingMutex.RLock()
	if _, isExist := p.receiving[in.Port]; !isExist {
		p.receivingMutex.RUnlock()

		storage := p.storage.List()
		uuid := uuid.NewV4().String()

		if in.StorageUUID != nil {
			uuid = in.StorageUUID.Value
			p.storagingMutex.RLock()
			if sTmp, isExist := p.storaging[uuid]; isExist {
				storage = sTmp
			}
			p.storagingMutex.RUnlock()
		}

		ctx, stop := context.WithCancel(context.Background())
		if err = p.udp.Receive(ctx, in.Port, storage); err == nil {
			p.storagingMutex.Lock()
			p.storaging[uuid] = storage
			p.storagingMutex.Unlock()

			p.receivingMutex.Lock()
			p.receiving[in.Port] = stop
			p.receivingMutex.Unlock()

			out = &StartReceiveResponse{
				StorageUUID: uuid,
			}
			return
		}
		stop()
		return
	}
	p.receivingMutex.RUnlock()
	err = fmt.Errorf("%v is busy", in.Port)
	return
}

// StopReceive stop receive data from server
func (p *player) StopReceive(c context.Context, in *StopReceiveRequest) (out *StopReceiveResponse, err error) {
	p.receivingMutex.RLock()
	if stop, isExist := p.receiving[in.Port]; isExist {
		p.receivingMutex.RUnlock()

		stop()

		p.receivingMutex.Lock()
		delete(p.receiving, in.Port)
		p.receivingMutex.Unlock()

		out = &StopReceiveResponse{}
		return
	}
	p.receivingMutex.RUnlock()
	err = fmt.Errorf("%v is not exist", in.Port)
	return
}

// StartPlay play audio on device
func (p *player) StartPlay(c context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	p.storagingMutex.RLock()
	storage, isExist := p.storaging[in.StorageUUID]
	p.storagingMutex.RUnlock()

	if !isExist {
		err = fmt.Errorf("storage %v is not exist", in.StorageUUID)
		return
	}

	p.playingMutex.RLock()
	if _, isExist := p.playing[in.DeviceName]; !isExist {
		p.playingMutex.RUnlock()

		ctx, stop := context.WithCancel(context.Background())
		if err = p.device.Play(ctx, in.DeviceName, int(in.Channels), int(in.Rate), storage); err == nil {
			p.playingMutex.Lock()
			p.playing[in.DeviceName] = stop
			p.playingMutex.Unlock()

			out = &StartPlayResponse{}
			return
		}
		stop()
		return
	}

	p.playingMutex.RUnlock()
	err = fmt.Errorf("%s is busy", in.DeviceName)
	return
}

// StopPlay stop play on device
func (p *player) StopPlay(c context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
	p.playingMutex.RLock()
	if stop, isExist := p.playing[in.DeviceName]; isExist {
		p.playingMutex.RUnlock()

		stop()

		p.playingMutex.Lock()
		delete(p.playing, in.DeviceName)
		p.playingMutex.Unlock()

		out = &StopPlayResponse{}
		return
	}
	p.playingMutex.RUnlock()
	err = fmt.Errorf("%s is not exist", in.DeviceName)
	return
}

// ClearStorage with StorageUUID
func (p *player) ClearStorage(c context.Context, in *ClearStorageRequest) (out *ClearStorageResponse, err error) {
	p.storagingMutex.RLock()
	if storage, isExist := p.storaging[in.StorageUUID]; isExist {
		p.storagingMutex.RUnlock()

		storage.Close()

		p.storagingMutex.Lock()
		delete(p.storaging, in.StorageUUID)
		p.storagingMutex.Unlock()

		out = &ClearStorageResponse{}
		return
	}
	p.storagingMutex.RUnlock()
	err = fmt.Errorf("%s is not exist", in.StorageUUID)
	return
}

// NewPlayer ...
func NewPlayer(
	udp udp,
	device device,
	storage storage,
) PlayerServer {
	return &player{
		receiving: make(map[string]context.CancelFunc),
		storaging: make(map[string]io.ReadWriteCloser),
		playing:   make(map[string]context.CancelFunc),

		udp:     udp,
		device:  device,
		storage: storage,
	}
}
