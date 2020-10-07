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
	// router.Handle(methodPlayerReceiveStop, uriPlayerReceiveStop, Handler(svc, newTransport(), ErrorProcessing))
	// router.Handle(methodPlayerPlay, uriPlayerPlay, Handler(svc, newTransport(), ErrorProcessing))
	// router.Handle(methodPlayerPause, uriPlayerPause, Handler(svc, newTransport(), ErrorProcessing))
	// router.Handle(methodPlayerStop, uriPlayerStop, Handler(svc, newTransport(), ErrorProcessing))

	// router.Handle(methodStartFileRecoding, uriStartFileRecoding, Handler(svc, newTransport(), ErrorProcessing))
	// router.Handle(methodStopFileRecoding, uriStopFileRecoding, Handler(svc, newTransport(), ErrorProcessing))
	// router.Handle(methodPlayFromRecorder, uriPlayFromRecorder, Handler(svc, newTransport(), ErrorProcessing))
	// router.Handle(methodStopFromRecorder, uriStopFromRecorder, Handler(svc, newTransport(), ErrorProcessing))

	// router.Handle(methodRecorderStart, uriRecorderStart, Handler(svc, newTransport(), ErrorProcessing))
	// router.Handle(methodRecorderStart, uriRecorderStart, Handler(svc, newTransport(), ErrorProcessing))

	router.Handle("GET", "/debug/pprof/", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index))
	router.Handle("GET", "/debug/pprof/profile", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile))

	return &fasthttp.Server{
		Handler: router.Handler,
	}
}
