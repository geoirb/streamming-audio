package player

import (
	"context"
	"fmt"
	"io"
	"sync"
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

// Player audio signal
type Player struct {
	mutex sync.Mutex
	port  map[string]context.CancelFunc

	udp     udp
	device  device
	storage storage
}

func (m *Player) StartRecieve(context.Context, *StartRecieveRequest) (*StartRecieveRequest, error) {}
func (m *Player) StopRecieve(context.Context, *StopRecieveRequest) (*StopRecieveRequest, error)    {}
func (m *Player) StartPlay(context.Context, *StartPlayRequest) (*StartPlayResponse, error)         {}
func (m *Player) StopPlay(context.Context, *StopPlayRequest) (*StopPlayResponse, error)            {}

// StartPlay play audio on device from server
func (m *Player) StartPlay(ctx context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, isExist := m.port[in.Port]; isExist {
		err = fmt.Errorf("%v is exist", in.Port)
		return
	}

	list := m.storage.List()
	c, cancel := context.WithCancel(context.Background())

	if err = m.udp.Receive(c, in.Port, list); err != nil {
		cancel()
		return
	}

	if err = m.device.Play(c, in.DeviceName, int(in.Channels), int(in.Rate), list); err != nil {
		cancel()
		return
	}

	m.port[in.Port] = cancel
	out = &StartPlayResponse{}
	return
}

// StopPlay stop play audio on device
func (m *Player) StopPlay(ctx context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	cancel, isExist := m.port[in.Port]
	if !isExist {
		err = fmt.Errorf("%v is exist", in.Port)
		return
	}
	cancel()
	delete(m.port, in.Port)
	out = &StopPlayResponse{}
	return
}

// NewPlayer ...
func NewPlayer(
	udp udp,
	device device,
	storage storage,
) PlayerServer {
	return &Player{
		port: make(map[string]context.CancelFunc),

		udp:     udp,
		device:  device,
		storage: storage,
	}
}
