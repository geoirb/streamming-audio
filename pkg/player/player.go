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

type receiving struct {
	isStart bool
	stop    context.CancelFunc
	storage io.ReadWriteCloser
}

// Player audio signal
type Player struct {
	receivingMutex sync.Mutex
	receiving      map[string]receiving
	playingMutex   sync.Mutex
	playing        map[string]context.CancelFunc

	udp     udp
	device  device
	storage storage
}

// StartReceive start receive data from server and save
func (p *Player) StartReceive(c context.Context, in *StartReceiveRequest) (out *StartReceiveRequest, err error) {
	p.receivingMutex.Lock()
	defer p.receivingMutex.Unlock()

	if _, isExist := p.receiving[in.Port]; !isExist {
		list := p.storage.List()
		ctx, cancel := context.WithCancel(context.Background())

		if err = p.udp.Receive(ctx, in.Port, list); err != nil {
			cancel()
			return
		}
		p.receiving[in.Port] = receiving{
			isStart: true,
			stop:    cancel,
			storage: list,
		}
		out = &StartReceiveRequest{}
		return
	}
	err = fmt.Errorf("%v is busy", in.Port)
	return
}

// StopReceive stop receive data from server
// Force in StopReceiveRequest delete storage
func (p *Player) StopReceive(c context.Context, in *StopReceiveRequest) (out *StopReceiveRequest, err error) {
	p.receivingMutex.Lock()
	defer p.receivingMutex.Unlock()

	if r, isExist := p.receiving[in.Port]; isExist {
		if r.isStart {
			r.stop()
		}
		if in.Force {
			delete(p.receiving, in.Port)
		}
		p.receiving[in.Port] = receiving{
			isStart: true,
			storage: r.storage,
		}
		out = &StopReceiveRequest{}
		return
	}
	err = fmt.Errorf("%v is not exist", in.Port)
	return
}

// StartPlay play audio on device
func (p *Player) StartPlay(c context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	p.receivingMutex.Lock()
	r, isExist := p.receiving[in.Port]
	if !isExist {
		p.receivingMutex.Unlock()
		err = fmt.Errorf("%v is not exist", in.Port)
		return
	}
	p.receivingMutex.Unlock()

	p.playingMutex.Lock()
	p.playingMutex.Unlock()

	if _, isExist := p.playing[in.DeviceName]; !isExist {
		ctx, cancel := context.WithCancel(context.Background())
		if err = p.device.Play(ctx, in.DeviceName, int(in.Channels), int(in.Rate), r.storage); err != nil {
			cancel()
			return
		}
		p.playing[in.DeviceName] = cancel
		out = &StartPlayResponse{}
		return
	}
	err = fmt.Errorf("%v is busy", in.Port)
	return
}

// StopPlay stop play on device
func (p *Player) StopPlay(c context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
	p.playingMutex.Lock()
	defer p.playingMutex.Unlock()

	if stop, isExist := p.playing[in.DeviceName]; isExist {
		stop()
		delete(p.playing, in.DeviceName)
		out = &StopPlayResponse{}
		return
	}
	err = fmt.Errorf("%v is not exist", in.DeviceName)
	return
}

// NewPlayer ...
func NewPlayer(
	udp udp,
	device device,
	storage storage,
) PlayerServer {
	return &Player{
		receiving: make(map[string]receiving),
		playing:   make(map[string]context.CancelFunc),

		udp:     udp,
		device:  device,
		storage: storage,
	}
}
