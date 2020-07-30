package media

import (
	"context"
	"time"
)

type connection interface {
	StartReceive() <-chan []byte
}

type device interface {
	Play(audio []int16)
}

type converter interface {
	ToInt16([]byte) []int16
}

type cash interface {
	Push(e []int16)
	Pop() []int16
}

// Media audio repicient
type Media struct {
	connection connection
	device     device
	converter  converter
	cash       cash
	size       int
}

// Repicenting data over vonnection
func (m *Media) Repicenting(ctx context.Context) {
	go func() {
		c := m.connection.StartReceive()
		for {
			select {
			case data := <-c:
				el := m.converter.ToInt16(data)
				m.cash.Push(el)
			case <-ctx.Done():
				return
			}
		}
	}()

	time.Sleep(20 * time.Second)
	for {
		audio := m.cash.Pop()
		if audio != nil {
			m.device.Play(audio)
		}
	}
}

// NewMedia ...
func NewMedia(
	connection connection,
	device device,
	converter converter,
	cash cash,
	size int,
) *Media {
	return &Media{
		connection: connection,
		device:     device,
		converter:  converter,
		cash:       cash,
		size:       size,
	}
}
