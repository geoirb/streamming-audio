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

type playerReceiveStop struct {
	svc             svc
	transport       PlayerReceiveStopTransport
	errorProcessing errorProcessing
}

func (s *playerReceiveStop) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                  error
		playerIP, playerPort string
	)
	if playerIP, playerPort, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.PlayerReceiveStop(ctx, playerIP, playerPort); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func playerReceiveStopHandler(svc svc, transport PlayerReceiveStopTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerReceiveStop{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerPlay struct {
	svc             svc
	transport       PlayerPlayTransport
	errorProcessing errorProcessing
}

func (s *playerPlay) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                              error
		playerIP, uuid, playerDeviceName string
		channels, rate                   uint32
	)
	if playerIP, uuid, playerDeviceName, channels, rate, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.PlayerPlay(ctx, playerIP, uuid, playerDeviceName, channels, rate); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func playerPlayHandler(svc svc, transport PlayerPlayTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerPlay{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerPause struct {
	svc             svc
	transport       PlayerPauseTransport
	errorProcessing errorProcessing
}

func (s *playerPause) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                        error
		playerIP, playerDeviceName string
	)
	if playerIP, playerDeviceName, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.PlayerPause(ctx, playerIP, playerDeviceName); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func playerPauseHandler(svc svc, transport PlayerPauseTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerPause{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerStop struct {
	svc             svc
	transport       PlayerStopTransport
	errorProcessing errorProcessing
}

func (s *playerStop) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                                          error
		playerIP, playerPort, playerDeviceName, uuid string
	)
	if playerIP, playerPort, playerDeviceName, uuid, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.PlayerStop(ctx, playerIP, playerPort, playerDeviceName, uuid); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func playerStopHandler(svc svc, transport PlayerStopTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerStop{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type startFileRecoding struct {
	svc             svc
	transport       StartFileRecodingTransport
	errorProcessing errorProcessing
}

func (s *startFileRecoding) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                                               error
		recorderIP, recorderDeviceName, receivePort, file string
		channels, rate                                    uint32
	)
	if recorderIP, recorderDeviceName, channels, rate, receivePort, file, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.StartFileRecoding(ctx, recorderIP, recorderDeviceName, channels, rate, receivePort, file); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func startFileRecodingHandler(svc svc, transport StartFileRecodingTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &startFileRecoding{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type stopFileRecoding struct {
	svc             svc
	transport       StopFileRecodingTransport
	errorProcessing errorProcessing
}

func (s *stopFileRecoding) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                                         error
		recorderIP, recorderDeviceName, receivePort string
	)
	if recorderIP, recorderDeviceName, receivePort, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.StopFileRecoding(ctx, recorderIP, recorderDeviceName, receivePort); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func stopFileRecodingHandler(svc svc, transport StopFileRecodingTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &stopFileRecoding{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playFromRecorder struct {
	svc             svc
	transport       PlayFromRecorderTransport
	errorProcessing errorProcessing
}

func (s *playFromRecorder) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                                                                          error
		playerIP, playerPort, playerDeviceName, recorderIP, recorderDeviceName, uuid string
		channels, rate                                                               uint32
	)
	if playerIP, playerPort, playerDeviceName, channels, rate, recorderIP, recorderDeviceName, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if uuid, err = s.svc.PlayFromRecorder(ctx, playerIP, playerPort, playerDeviceName, channels, rate, recorderIP, recorderDeviceName); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response, uuid); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func playFromRecorderHandler(svc svc, transport PlayFromRecorderTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playFromRecorder{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type stopFromRecorder struct {
	svc             svc
	transport       StopFromRecorderTransport
	errorProcessing errorProcessing
}

func (s *stopFromRecorder) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                                                                          error
		playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string
	)
	if playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.StopFromRecorder(ctx, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func stopFromRecorderHandler(svc svc, transport StopFromRecorderTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &stopFromRecorder{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type recorderStart struct {
	svc             svc
	transport       RecorderStartTransport
	errorProcessing errorProcessing
}

func (s *recorderStart) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                                     error
		recorderIP, recorderDeviceName, dstAddr string
		channels, rate                          uint32
	)
	if recorderIP, recorderDeviceName, channels, rate, dstAddr, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.RecorderStart(ctx, recorderIP, recorderDeviceName, channels, rate, dstAddr); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func recorderStartHandler(svc svc, transport RecorderStartTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &recorderStart{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type recorderStop struct {
	svc             svc
	transport       RecorderStopTransport
	errorProcessing errorProcessing
}

func (s *recorderStop) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                            error
		recorderIP, recorderDeviceName string
	)
	if recorderIP, recorderDeviceName, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.RecoderStop(ctx, recorderIP, recorderDeviceName); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func recorderStopHandler(svc svc, transport RecorderStopTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &recorderStop{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}
