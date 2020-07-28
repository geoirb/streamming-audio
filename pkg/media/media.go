package media

import (
	"context"
	"fmt"
)

type converter interface {
	ToInt16(src []byte) (dst []int16)
}

type device interface {
	Play(audio []int16)
}

type connection interface {
	StartReceive() <-chan []byte
}

// Media audio repicient
type Media struct {
	connection connection
	device     device
	converter  converter
	size       int
}

// Repicenting data over vonnection
func (m *Media) Repicenting(ctx context.Context) (err error) {
	c := m.connection.StartReceive()
	for {
		select {
		case data := <-c:
			fmt.Println(data)
			audio := m.converter.ToInt16(data)
			m.device.Play(audio)
		case <-ctx.Done():
			return
		}
	}
}

// NewMedia ...
func NewMedia(
	connection connection,
	device device,
	converter converter,

	size int,
) *Media {
	return &Media{
		connection: connection,
		device:     device,
		converter:  converter,
		size:       size,
	}
}
