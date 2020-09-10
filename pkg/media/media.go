package media

import (
	"context"
	"fmt"
	"sync"

	"github.com/geoirb/sound-ethernet-streaming/pkg/cash"
)

type client interface {
	Receive(context.Context, string) (<-chan []byte, error)
}

type device interface {
	Write([]byte)
}

// Media audio receiver
type Media struct {
	mutex   sync.Mutex
	receive map[string]context.CancelFunc
	client  client
}

// Add ...
func (m *Media) Add(ctx context.Context, port string) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, isExist := m.receive[port]; isExist {
		err = fmt.Errorf("receive port is exist: %v", port)
		return
	}

	c, cancel := context.WithCancel(ctx)
	data, err := m.client.Receive(c, port)
	if err != nil {
		cancel()
		return
	}

	c2h := cash.NewCash()

	go m.receiving(c, data, c2h)
	m.receive[port] = cancel
	return nil
}

func (m *Media) receiving(ctx context.Context, data <-chan []byte, c2h *cash.Cash) {
	for {
		select {
		case sample := <-data:
			c2h.Push(sample)
		case <-ctx.Done():
			return
		}
	}
}

// func (m *Media) play(ctx context.Context, device device, cash cash) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			if samples := cash.Pop(); samples != nil && len(samples) > 0 {
// 				device.Write(samples)
// 			}
// 		}
// 	}
// }

// NewClient ...
func NewClient() *Media {
	return &Media{
		receive: make(map[string]context.CancelFunc),
	}
}
