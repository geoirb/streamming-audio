package httpserver

import (
	"net/http"

	"github.com/valyala/fasthttp"

	"github.com/geoirb/ausio-service/pkg/server"
)

type filePlay struct {
	svc             server.Server
	transport       FilePlayTransport
	errorProcessing errorProcessing
}

func (s *filePlay) handler(ctx *fasthttp.RequestCtx) {
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

	if uuid, channels, rate, err = s.svc.FilePlay(ctx, file, playerIP, playerPort, playerDeviceName); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response, uuid, channels, rate); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func filePlayHandler(svc server.Server, transport FilePlayTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &filePlay{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type fileStop struct {
	svc             server.Server
	transport       FileStopTransport
	errorProcessing errorProcessing
}

func (s *fileStop) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                                          error
		playerIP, playerPort, playerDeviceName, uuid string
	)
	if playerIP, playerPort, playerDeviceName, uuid, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.FileStop(ctx, playerIP, playerPort, playerDeviceName, uuid); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func fileStopHandler(svc server.Server, transport FileStopTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &fileStop{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerReceiveStart struct {
	svc             server.Server
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

func playerReceiveStartHandler(svc server.Server, transport PlayerReceiveStartTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerReceiveStart{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerReceiveStop struct {
	svc             server.Server
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

func playerReceiveStopHandler(svc server.Server, transport PlayerReceiveStopTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerReceiveStop{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerPlay struct {
	svc             server.Server
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

func playerPlayHandler(svc server.Server, transport PlayerPlayTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerPlay{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerStop struct {
	svc             server.Server
	transport       PlayerStopTransport
	errorProcessing errorProcessing
}

func (s *playerStop) handler(ctx *fasthttp.RequestCtx) {
	var (
		err                        error
		playerIP, playerDeviceName string
	)
	if playerIP, playerDeviceName, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.PlayerStop(ctx, playerIP, playerDeviceName); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func playerStopHandler(svc server.Server, transport PlayerStopTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerStop{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playerClearStorage struct {
	svc             server.Server
	transport       PlayerClearStorageTransport
	errorProcessing errorProcessing
}

func (s *playerClearStorage) handler(ctx *fasthttp.RequestCtx) {
	var (
		err            error
		playerIP, uuid string
	)
	if playerIP, uuid, err = s.transport.Decode(ctx); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusBadRequest)
		return
	}

	if err = s.svc.PlayerClearStorage(ctx, playerIP, uuid); err != nil {
		s.errorProcessing(&ctx.Response, err, -1)
		return
	}

	if err = s.transport.Encode(&ctx.Response); err != nil {
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

func playerClearStorageHandler(svc server.Server, transport PlayerClearStorageTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playerClearStorage{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type startFileRecoding struct {
	svc             server.Server
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

func startFileRecodingHandler(svc server.Server, transport StartFileRecodingTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &startFileRecoding{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type stopFileRecoding struct {
	svc             server.Server
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

func stopFileRecodingHandler(svc server.Server, transport StopFileRecodingTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &stopFileRecoding{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type playFromRecorder struct {
	svc             server.Server
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

func playFromRecorderHandler(svc server.Server, transport PlayFromRecorderTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &playFromRecorder{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type stopFromRecorder struct {
	svc             server.Server
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

func stopFromRecorderHandler(svc server.Server, transport StopFromRecorderTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &stopFromRecorder{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type recorderStart struct {
	svc             server.Server
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

func recorderStartHandler(svc server.Server, transport RecorderStartTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &recorderStart{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}

type recorderStop struct {
	svc             server.Server
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

func recorderStopHandler(svc server.Server, transport RecorderStopTransport, errorProcessing errorProcessing) fasthttp.RequestHandler {
	s := &recorderStop{
		svc:             svc,
		transport:       transport,
		errorProcessing: errorProcessing,
	}
	return s.handler
}
