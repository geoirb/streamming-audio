package media

import (
	"context"
	"fmt"
	"sync"

	"github.com/geoirb/sound-ethernet-streaming/pkg/cash"
	"github.com/geoirb/sound-ethernet-streaming/pkg/controller/media"
)

type client interface {
	Receive(context.Context, string, *cash.Cash) error
}

type device interface {
	Play(context.Context, string, int, int, *cash.Cash) error
}

// Media audio receiver
type Media struct {
	mutex   sync.Mutex
	receive map[string]context.CancelFunc
	client  client
	device  device
}

// StartReceive audio data from port
func (m *Media) StartReceive(ctx context.Context, in *media.StartReceiveRequest) (out *media.StartReceiveResponse, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, isExist := m.receive[in.Port]; isExist {
		err = fmt.Errorf("receive port is exist: %v", in.Port)
		return
	}

	c2h := cash.NewCash()
	c, cancel := context.WithCancel(ctx)

	err = m.client.Receive(c, in.Port, c2h)
	if err != nil {
		cancel()
		return
	}

	err = m.device.Play(c, in.DeviceName, int(in.Channels), int(in.Rate), c2h)
	if err != nil {
		cancel()
		return
	}

	m.receive[in.Port] = cancel
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
) *Media {
	return &Media{
		receive: make(map[string]context.CancelFunc),
		client:  client,
		device:  device,
	}
}
