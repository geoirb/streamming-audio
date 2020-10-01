package recorder

import (
	"context"
	"fmt"
	"io"
	"sync"
)

type udp interface {
	TurnOnSender(string) (io.WriteCloser, error)
}

type device interface {
	Record(context.Context, string, int, int, io.WriteCloser) error
}

// Recorder audio signal
type Recorder struct {
	mutex    sync.Mutex
	recoding map[string]context.CancelFunc

	udp    udp
	device device
}

// StartRecord ...
func (r *Recorder) StartRecord(c context.Context, in *StartRecordRequest) (out *StartRecordResponse, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, isExist := r.recoding[in.DeviceName]; !isExist {
		var destination io.WriteCloser
		if destination, err = r.udp.TurnOnSender(in.DestAddr); err == nil {
			ctx, stop := context.WithCancel(context.Background())
			if err = r.device.Record(ctx, in.DeviceName, int(in.Channels), int(in.Rate), destination); err == nil {
				r.recoding[in.DeviceName] = stop
				out = &StartRecordResponse{}
				return
			}
			stop()
		}
		return
	}
	err = fmt.Errorf("%v is busy", in.DeviceName)
	return
}

// StopRecord ...
func (r *Recorder) StopRecord(ctx context.Context, in *StopRecordRequest) (out *StopRecordResponse, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if stop, isExist := r.recoding[in.DeviceName]; isExist {
		stop()
		delete(r.recoding, in.DeviceName)
		out = &StopRecordResponse{}
		return
	}
	err = fmt.Errorf("%v is not exist", in.DeviceName)
	return
}

// NewRecorder ...
func NewRecorder(
	udp udp,
	device device,
) RecorderServer {
	return &Recorder{
		recoding: make(map[string]context.CancelFunc),

		udp:    udp,
		device: device,
	}
}
