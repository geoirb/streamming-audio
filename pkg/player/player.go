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
	Receive(ctx context.Context, receivePort string, storage io.Writer) error
}

type device interface {
	Play(context.Context, string, int, int, io.Reader) error
}

type player struct {
	receivingMutex sync.Mutex
	receiving      map[string]func()

	storagingMutex sync.Mutex
	storaging      map[string]io.ReadWriteCloser

	playingMutex sync.Mutex
	playing      map[string]func()

	udp     udp
	device  device
	storage storage
}

// ReceiveStart start receive data from server and save
func (p *player) ReceiveStart(c context.Context, in *StartReceiveRequest) (out *StartReceiveResponse, err error) {
	p.receivingMutex.Lock()
	defer p.receivingMutex.Unlock()

	if _, isExist := p.receiving[in.Port]; !isExist {
		storage := p.storage.List()
		uuid := uuid.NewV4().String()

		if in.StorageUUID != nil {
			uuid = in.StorageUUID.Value
			if sTmp, isExist := p.storaging[uuid]; isExist {
				storage = sTmp
			}
		}

		ctx, stop := context.WithCancel(context.Background())
		if err = p.udp.Receive(ctx, in.Port, storage); err == nil {
			p.storaging[uuid] = storage
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

// ReceiveStop stop receive data from server
func (p *player) ReceiveStop(c context.Context, in *StopReceiveRequest) (out *StopReceiveResponse, err error) {
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

// Play play audio on device
func (p *player) Play(c context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	p.storagingMutex.Lock()
	defer p.storagingMutex.Unlock()

	storage, isExist := p.storaging[in.StorageUUID]
	if !isExist {
		err = fmt.Errorf("storage %v is not exist", in.StorageUUID)
		return
	}

	if _, isExist := p.playing[in.DeviceName]; !isExist {
		ctx, stop := context.WithCancel(context.Background())
		if err = p.device.Play(ctx, in.DeviceName, int(in.Channels), int(in.Rate), storage); err == nil {
			p.playing[in.DeviceName] = stop
			out = &StartPlayResponse{}
			return
		}
		stop()
		return
	}
	err = fmt.Errorf("%s is busy", in.DeviceName)
	return
}

// Stop stop play on device
func (p *player) Stop(c context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
	p.playingMutex.Lock()
	defer p.playingMutex.Unlock()

	if stop, isExist := p.playing[in.DeviceName]; isExist {
		stop()
		delete(p.playing, in.DeviceName)
		out = &StopPlayResponse{}
		return
	}
	err = fmt.Errorf("%s is not exist", in.DeviceName)
	return
}

// ClearStorage with StorageUUID
func (p *player) ClearStorage(c context.Context, in *ClearStorageRequest) (out *ClearStorageResponse, err error) {
	p.storagingMutex.Lock()
	defer p.storagingMutex.Unlock()

	if storage, isExist := p.storaging[in.StorageUUID]; isExist {
		storage.Close()
		delete(p.storaging, in.StorageUUID)
		out = &ClearStorageResponse{}
		return
	}
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
		receiving: make(map[string]func()),
		storaging: make(map[string]io.ReadWriteCloser),
		playing:   make(map[string]func()),

		udp:     udp,
		device:  device,
		storage: storage,
	}
}
