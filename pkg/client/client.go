package client

import (
	"context"
	"fmt"
)

type connection interface {
	StartReceive(context.Context)
	Data() <-chan []byte
}

type device interface {
	Write([]byte)
}

type cash interface {
	Push([]byte)
	Pop() []byte
}

type receive struct {
	connection connection
	cash       cash
}

// Client audio receiver
type Client struct {
	pull map[device]receive
}

// Add ...
func (m *Client) Add(device device, connection connection, cash cash) error {
	if _, isExist := m.pull[device]; isExist {
		return fmt.Errorf("device is exist: %v", connection)
	}
	m.pull[device] = receive{
		connection: connection,
		cash:       cash,
	}
	return nil
}

// Start client
func (m *Client) Start(ctx context.Context) {
	for device, i := range m.pull {
		go m.receiving(ctx, i.connection, i.cash)
		go m.play(ctx, device, i.cash)
	}
}

func (m *Client) receiving(ctx context.Context, connection connection, cash cash) {
	go connection.StartReceive(ctx)
	for {
		select {
		case data := <-connection.Data():
			cash.Push(data)
		case <-ctx.Done():
			return
		}
	}
}

func (m *Client) play(ctx context.Context, device device, cash cash) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if samples := cash.Pop(); samples != nil && len(samples) > 0 {
				device.Write(samples)
			}
		}
	}
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		pull: make(map[device]receive),
	}
}
