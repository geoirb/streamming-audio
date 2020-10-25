package player

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/twinj/uuid"
)

type storageCreator interface {
	List() io.ReadWriteCloser
}

type tcp interface {
	Receive(ctx context.Context, receivePort string, storage io.Writer) error
}

type device interface {
	Play(ctx context.Context, deviceName string, channels, rate, bitsPerSample int, r io.Reader) (err error)
}

type player struct {
	receivingMutex sync.Mutex
	receivingPort  map[string]func()

	storageMutex sync.Mutex
	storage      map[string]io.ReadWriteCloser

	playbackDeviceMutex sync.Mutex
	playbackDevice      map[string]func()

	tcp            tcp
	device         device
	storageCreator storageCreator
}

// State return all busy ports, devices on player and existing storage
func (p *player) State(ctx context.Context, in *StateRequest) (out *StateResponse, err error) {
	out = &StateResponse{}

	p.receivingMutex.Lock()
	out.Ports = make([]string, 0, len(p.receivingPort))
	for port := range p.receivingPort {
		out.Ports = append(out.Ports, port)
	}
	p.receivingMutex.Unlock()

	p.storageMutex.Lock()
	out.Storages = make([]string, 0, len(p.storage))
	for uuid := range p.receivingPort {
		out.Storages = append(out.Storages, uuid)
	}
	p.storageMutex.Unlock()

	p.playbackDeviceMutex.Lock()
	out.Devices = make([]string, 0, len(p.playbackDevice))
	for device := range p.playbackDevice {
		out.Devices = append(out.Devices, device)
	}
	p.playbackDeviceMutex.Unlock()
	return
}

// ReceiveStart start receive data from server and save
func (p *player) ReceiveStart(c context.Context, in *StartReceiveRequest) (out *StartReceiveResponse, err error) {
	p.receivingMutex.Lock()
	defer p.receivingMutex.Unlock()

	if _, isExist := p.receivingPort[in.Port]; !isExist {
		storage := p.storageCreator.List()
		uuid := uuid.NewV4().String()

		if in.StorageUUID != nil {
			uuid = in.StorageUUID.Value
			if sTmp, isExist := p.storage[uuid]; isExist {
				storage = sTmp
			}
		}

		ctx, stop := context.WithCancel(context.Background())
		if err = p.tcp.Receive(ctx, in.Port, storage); err == nil {
			p.storage[uuid] = storage
			p.receivingPort[in.Port] = stop
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

	if stop, isExist := p.receivingPort[in.Port]; isExist {
		stop()
		delete(p.receivingPort, in.Port)
		out = &StopReceiveResponse{}
		return
	}
	err = fmt.Errorf("%v is not exist", in.Port)
	return
}

// Play play audio on device
func (p *player) Play(c context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	p.storageMutex.Lock()
	defer p.storageMutex.Unlock()

	storage, isExist := p.storage[in.StorageUUID]
	if !isExist {
		err = fmt.Errorf("storage %v is not exist", in.StorageUUID)
		return
	}

	if _, isExist := p.playbackDevice[in.DeviceName]; !isExist {
		ctx, stop := context.WithCancel(context.Background())
		if err = p.device.Play(ctx, in.DeviceName, int(in.Channels), int(in.Rate), int(in.BitsPerSample), storage); err == nil {
			p.playbackDevice[in.DeviceName] = stop
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
	p.playbackDeviceMutex.Lock()
	defer p.playbackDeviceMutex.Unlock()

	if stop, isExist := p.playbackDevice[in.DeviceName]; isExist {
		stop()
		delete(p.playbackDevice, in.DeviceName)
		out = &StopPlayResponse{}
		return
	}
	err = fmt.Errorf("%s is not exist", in.DeviceName)
	return
}

// ClearStorage with StorageUUID
func (p *player) ClearStorage(c context.Context, in *ClearStorageRequest) (out *ClearStorageResponse, err error) {
	p.storageMutex.Lock()
	defer p.storageMutex.Unlock()

	if storage, isExist := p.storage[in.StorageUUID]; isExist {
		storage.Close()
		delete(p.storage, in.StorageUUID)
		out = &ClearStorageResponse{}
		return
	}
	err = fmt.Errorf("%s is not exist", in.StorageUUID)
	return
}

// NewPlayer ...
func NewPlayer(
	tcp tcp,
	device device,
	storage storageCreator,
) PlayerServer {
	return &player{
		receivingPort:  make(map[string]func()),
		storage:        make(map[string]io.ReadWriteCloser),
		playbackDevice: make(map[string]func()),

		tcp:            tcp,
		device:         device,
		storageCreator: storage,
	}
}
