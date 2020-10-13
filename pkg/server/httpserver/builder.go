package httpserver

import (
	"net/http"
	"net/http/pprof"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/geoirb/sound-server/pkg/server"
)

const (
	methodFilePlay = http.MethodPost
	uriFilePlay    = "/player/file/play"
	methodFileStop = http.MethodPost
	uriFileStop    = "/player/file/stop"

	methodPlayerReceiveStart = http.MethodPost
	uriPlayerReceiveStart    = "/player/receive/start"
	methodPlayerReceiveStop  = http.MethodPost
	uriPlayerReceiveStop     = "/player/receive/stop"
	methodPlayerPlay         = http.MethodPost
	uriPlayerPlay            = "/player/play"
	methodPlayerStop         = http.MethodPost
	uriPlayerStop            = "/player/stop"
	methodPlayerClearStorage = http.MethodPost
	uriPlayerClearStorage    = "/player/clearstorage"

	methodStartFileRecoding = http.MethodPost
	uriStartFileRecoding    = "/recoder/file/start"
	methodStopFileRecoding  = http.MethodPost
	uriStopFileRecoding     = "/recoder/file/stop"
	methodPlayFromRecorder  = http.MethodPost
	uriPlayFromRecorder     = "/recoder/player/play"
	methodStopFromRecorder  = http.MethodPost
	uriStopFromRecorder     = "/recoder/player/stop"

	methodRecorderStart = http.MethodPost
	uriRecorderStart    = "/recoder/start"
	methodRecoderStop   = http.MethodPost
	uriRecorderStop     = "/recoder/stop"
)

// NewServer return http server
func NewServer(svc server.Server) *fasthttp.Server {
	router := fasthttprouter.New()

	router.Handle(methodFilePlay, uriFilePlay, filePlayHandler(svc, newFilePlayTransport(), ErrorProcessing))
	router.Handle(methodFileStop, uriFileStop, fileStopHandler(svc, newFileStopTransport(), ErrorProcessing))

	router.Handle(methodPlayerReceiveStart, uriPlayerReceiveStart, playerReceiveStartHandler(svc, newPlayerReceiveStartTransport(), ErrorProcessing))
	router.Handle(methodPlayerReceiveStop, uriPlayerReceiveStop, playerReceiveStopHandler(svc, newPlayerReceiveStopTransport(), ErrorProcessing))
	router.Handle(methodPlayerPlay, uriPlayerPlay, playerPlayHandler(svc, newPlayerPlayTransport(), ErrorProcessing))
	router.Handle(methodPlayerStop, uriPlayerStop, playerStopHandler(svc, newPlayerStopTransport(), ErrorProcessing))
	router.Handle(methodPlayerClearStorage, uriPlayerClearStorage, playerClearStorageHandler(svc, newPlayerClearStorageTransport(), ErrorProcessing))

	router.Handle(methodStartFileRecoding, uriStartFileRecoding, startFileRecodingHandler(svc, newStartFileRecodingTransport(), ErrorProcessing))
	router.Handle(methodStopFileRecoding, uriStopFileRecoding, stopFileRecodingHandler(svc, newStopFileRecodingTransport(), ErrorProcessing))
	router.Handle(methodPlayFromRecorder, uriPlayFromRecorder, playFromRecorderHandler(svc, newPlayFromRecorderTransport(), ErrorProcessing))
	router.Handle(methodStopFromRecorder, uriStopFromRecorder, stopFromRecorderHandler(svc, newStopFromRecorderTransport(), ErrorProcessing))

	router.Handle(methodRecorderStart, uriRecorderStart, recorderStartHandler(svc, newRecorderStartTransport(), ErrorProcessing))
	router.Handle(methodRecoderStop, uriRecorderStop, recorderStopHandler(svc, newRecorderStopTransport(), ErrorProcessing))

	router.Handle("GET", "/debug/pprof/", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index))
	router.Handle("GET", "/debug/pprof/profile", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile))

	return &fasthttp.Server{
		Handler: router.Handler,
	}
}
