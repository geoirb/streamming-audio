package recoder

import (
	"context"
	"io"
	"sync"

	"github.com/geoirb/sound-ethernet-streaming/pkg/recoder/grpc"
)

type storage interface {
	List() io.ReadWriteCloser
}

type udp interface {
	Send(context.Context, string, io.Reader) error
}

type device interface {
	Recode(context.Context, string, int, int, io.Writer) error
}

// Recoder audio signal
type Recoder struct {
	mutex  sync.Mutex
	server map[string]context.CancelFunc

	device device
}

// StartRecode ...
func (r *Recoder) StartRecode(ctx context.Context, in *grpc.StartRecodeRequest) (out *grpc.StartRecodeResponse, err error) {
	return
}

// StopRecode ...
func (r *Recoder) StopRecode(ctx context.Context, in *grpc.StopRecodeRequest) (out *grpc.StopRecodeResponse, err error) {
	return
}
