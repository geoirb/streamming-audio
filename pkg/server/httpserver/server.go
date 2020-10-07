package httpserver

import (
	"context"
	"net/http"

	"github.com/valyala/fasthttp"
)

type svc interface {
	FilePlaying(ctx context.Context, file, playerIP, playerPort, playerDeviceName string) (uuid string, channels uint16, rate uint32, err error)

	PlayerReceiveStart(ctx context.Context, playerIP, playerPort string, uuid *string) (string, error)
	PlayerReceiveStop(ctx context.Context, playerIP, playerPort string) error
	PlayerPlay(ctx context.Context, playerIP, uuid, playerDeviceName string, channels, rate uint32) (err error)
	PlayerPause(ctx context.Context, playerIP, playerDeviceName string) (err error)
	PlayerStop(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid string) (err error)

	StartFileRecoding(ctx context.Context, recorderIP, recorderDeviceName string, channels, rate uint32, receivePort, file string) (err error)
	StopFileRecoding(ctx context.Context, recorderIP, recorderDeviceName, receivePort string) error
	PlayFromRecorder(ctx context.Context, playerIP, playerPort, playerDeviceName string, channels, rate uint32, recorderIP, recorderDeviceName string) (uuid string, err error)
	StopFromRecorder(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string) error

	RecorderStart(ctx context.Context, recorderIP, recorderDeviceName string, channels, rate uint32, dstAddr string) error
	RecoderStop(ctx context.Context, recorderIP, recorderDeviceName string) error
}

type filePlaying struct {
	svc             svc
	transport       FilePlayingTransport
	errorProcessing errorProcessing
}

func (s *filePlaying) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                                                error
		file, playerIP, playerPort, playerDeviceName, uuid string
		channels                                           uint16
		rate                                               uint32
	)
	if file, playerIP, playerPort, playerDeviceName, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if uuid, channels, rate, err = s.svc.FilePlaying(ctx, file, playerIP, playerPort, playerDeviceName); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response, uuid, channels, rate); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func filePlayingHandler(svc svc, transport FilePlayingTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &filePlaying{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerReceiveStart struct {
	svc             svc
	transport       PlayerReceiveStartTransport
	errorProcessing errorProcessing
}

func (s *playerReceiveStart) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                         error
		playerIP, playerPort, sUUID string
		uuid                        *string
	)
	if playerIP, playerPort, uuid, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if sUUID, err = s.svc.PlayerReceiveStart(ctx, playerIP, playerPort, uuid); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response, sUUID); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func playerReceiveStartHandler(svc svc, transport PlayerReceiveStartTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerReceiveStart{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}
