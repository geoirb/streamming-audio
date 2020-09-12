package media

import (
	"context"
	"fmt"
	"sync"

	"github.com/geoirb/sound-ethernet-streaming/pkg/controller/media"
	s "github.com/geoirb/sound-ethernet-streaming/pkg/storage"
)

type storage interface {
	NewList() s.List
}

type client interface {
	Receive(context.Context, string, s.List) error
}

type device interface {
	Play(context.Context, string, int, int, s.List) error
}

// Media audio receiver
type Media struct {
	mutex   sync.Mutex
	receive map[string]context.CancelFunc
	client  client
	device  device
	storage storage
}

// StartReceive audio data from port
func (m *Media) StartReceive(ctx context.Context, in *media.StartReceiveRequest) (out *media.StartReceiveResponse, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, isExist := m.receive[in.Port]; isExist {
		err = fmt.Errorf("receive port is exist: %v", in.Port)
		return
	}

	list := m.storage.NewList()
	c, cancel := context.WithCancel(context.Background())

	err = m.client.Receive(c, in.Port, list)
	if err != nil {
		cancel()
		return
	}

	err = m.device.Play(c, in.DeviceName, int(in.Channels), int(in.Rate), list)
	if err != nil {
		cancel()
		return
	}

	m.receive[in.Port] = cancel
	out = &media.StartReceiveResponse{}
	return
}

// StopReceive from port
func (m *Media) StopReceive(ctx context.Context, in *media.StopReceiveRequest) (out *media.StopReceiveResponse, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	cancel, isExist := m.receive[in.Port]
	if isExist {
		err = fmt.Errorf("receive port is exist: %v", in.Port)
		return
	}
	cancel()
	delete(m.receive, in.Port)
	return
}

// NewMedia ...
func NewMedia(
	client client,
	device device,
	storage storage,
) *Media {
	return &Media{
		receive: make(map[string]context.CancelFunc),
		client:  client,
		device:  device,
		storage: storage,
	}
}
