package httpserver

import (
	"net/http"
	"net/http/pprof"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

const (
	methodFilePlaying = http.MethodPost
	uriFilePlaying    = "/player/file"

	methodPlayerReceiveStart = http.MethodPost
	uriPlayerReceiveStart    = "/player/receive/start"
	methodPlayerReceiveStop  = http.MethodPost
	uriPlayerReceiveStop     = "/player/receive/stop"
	methodPlayerPlay         = http.MethodPost
	uriPlayerPlay            = "/player/play"
	methodPlayerPause        = http.MethodPost
	uriPlayerPause           = "/player/pause"
	methodPlayerStop         = http.MethodPost
	uriPlayerStop            = "/player/stop"

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
func NewServer(svc svc) *fasthttp.Server {
	router := fasthttprouter.New()

	router.Handle(methodFilePlaying, uriFilePlaying, filePlayingHandler(svc, newFilePlayingTransport(), ErrorProcessing))

	router.Handle(methodPlayerReceiveStart, uriPlayerReceiveStart, playerReceiveStartHandler(svc, newPlayerReceiveStartTransport(), ErrorProcessing))
	router.Handle(methodPlayerReceiveStop, uriPlayerReceiveStop, playerReceiveStopHandler(svc, newPlayerReceiveStopTransport(), ErrorProcessing))
	router.Handle(methodPlayerPlay, uriPlayerPlay, playerPlayHandler(svc, newPlayerPlayTransport(), ErrorProcessing))
	router.Handle(methodPlayerPause, uriPlayerPause, playerPauseHandler(svc, newPlayerPauseTransport(), ErrorProcessing))
	router.Handle(methodPlayerStop, uriPlayerStop, playerStopHandler(svc, newPlayerStopTransport(), ErrorProcessing))

	router.Handle(methodStartFileRecoding, uriStartFileRecoding, startFileRecodingHandler(svc, newStartFileRecodingTransport(), ErrorProcessing))
	router.Handle(methodStopFileRecoding, uriStopFileRecoding, stopFileRecodingHandler(svc, newStopFileRecodingTransport(), ErrorProcessing))
	router.Handle(methodStopFileRecoding, uriStopFileRecoding, playFromRecorderHandler(svc, newPlayFromRecorderTransport(), ErrorProcessing))
	router.Handle(methodStopFileRecoding, uriStopFileRecoding, stopFromRecorderHandler(svc, newStopFromRecorderTransport(), ErrorProcessing))

	router.Handle(methodRecorderStart, uriRecorderStart, recorderStartHandler(svc, newRecorderStartTransport(), ErrorProcessing))
	router.Handle(methodRecorderStart, uriRecorderStart, recorderStopHandler(svc, newRecorderStopTransport(), ErrorProcessing))

	router.Handle("GET", "/debug/pprof/", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index))
	router.Handle("GET", "/debug/pprof/profile", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile))

	return &fasthttp.Server{
		Handler: router.Handler,
	}
}
