package recoder

import (
	"context"
	"fmt"
	"io"
	"sync"
)

type udp interface {
	TurnOnSender(string) (io.ReadWriteCloser, error)
}

type device interface {
	Recode(context.Context, string, int, int, io.ReadWriteCloser) error
}

// Recoder audio signal
type Recoder struct {
	mutex  sync.Mutex
	server map[string]context.CancelFunc

	udp    udp
	device device
}

// StartRecode ...
func (r *Recoder) StartRecode(ctx context.Context, in *StartRecodeRequest) (out *StartRecodeResponse, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, isExist := r.server[in.DestAddr]; isExist {
		err = fmt.Errorf("%v is exist", in.DestAddr)
		return
	}

	conn, err := r.udp.TurnOnSender(in.DestAddr)
	if err != nil {
		return
	}

	c, cancel := context.WithCancel(context.Background())
	if err = r.device.Recode(c, in.DeviceName, int(in.Channels), int(in.Rate), conn); err != nil {
		cancel()
		return
	}

	r.server[in.DestAddr] = cancel
	out = &StartRecodeResponse{}
	return
}

// StopRecode ...
func (r *Recoder) StopRecode(ctx context.Context, in *StopRecodeRequest) (out *StopRecodeResponse, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	cancel, isExist := r.server[in.DestAddr]
	if !isExist {
		err = fmt.Errorf("%v is not exist", in.DestAddr)
		return
	}
	cancel()
	delete(r.server, in.DestAddr)
	out = &StopRecodeResponse{}
	return
}

// NewRecoder ...
func NewRecoder(
	udp udp,
	device device,
) RecoderServer {
	return &Recoder{
		server: make(map[string]context.CancelFunc),

		udp:    udp,
		device: device,
	}
}
