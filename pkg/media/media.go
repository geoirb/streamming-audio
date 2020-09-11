package media

import (
	"context"
	"fmt"
	"sync"

	"github.com/geoirb/sound-ethernet-streaming/pkg/cash"
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

// Receive audio data from port 
func (m *Media) Receive(ctx context.Context, port string, deviceName string, channels, rate int) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, isExist := m.receive[port]; isExist {
		err = fmt.Errorf("receive port is exist: %v", port)
		return
	}

	c2h := cash.NewCash()
	c, cancel := context.WithCancel(ctx)

	err = m.client.Receive(c, port, c2h)
	if err != nil {
		cancel()
		return
	}

	err = m.device.Play(c, deviceName, channels, rate, c2h)
	if err != nil {
		cancel()
		return
	}

	m.receive[port] = cancel
	return nil
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
