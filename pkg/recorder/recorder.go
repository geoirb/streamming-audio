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

type recorder struct {
	mutex         sync.Mutex
	captureDevice map[string]func()

	udp    udp
	device device
}

// State return busy recorder device
func (r *recorder) State(ctx context.Context, in *StateRequest) (out *StateResponse, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	out = &StateResponse{
		Devices: make([]string, 0, len(r.captureDevice)),
	}
	for device := range r.captureDevice {
		out.Devices = append(out.Devices, device)
	}
	return
}

// Start recording audio on recorder from recorderDeviceName
func (r *recorder) Start(c context.Context, in *StartSendRequest) (out *StartSendResponse, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, isExist := r.captureDevice[in.DeviceName]; !isExist {
		var destination io.WriteCloser
		if destination, err = r.udp.TurnOnSender(in.DestAddr); err == nil {
			ctx, stop := context.WithCancel(context.Background())
			if err = r.device.Record(ctx, in.DeviceName, int(in.Channels), int(in.Rate), destination); err == nil {
				r.captureDevice[in.DeviceName] = stop
				out = &StartSendResponse{}
				return
			}
			stop()
		}
		return
	}
	err = fmt.Errorf("%v is busy", in.DeviceName)
	return
}

// Stop recording audio on recorder from recorderDeviceName
func (r *recorder) Stop(ctx context.Context, in *StopSendRequest) (out *StopSendResponse, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if stop, isExist := r.captureDevice[in.DeviceName]; isExist {
		stop()
		delete(r.captureDevice, in.DeviceName)
		out = &StopSendResponse{}
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
	return &recorder{
		captureDevice: make(map[string]func()),

		udp:    udp,
		device: device,
	}
}
