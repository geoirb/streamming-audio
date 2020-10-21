package httpclient

import (
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

const (
	protocol = "http"

	methodFilePlay = http.MethodPost
	uriFilePlay    = "/player/file/play"
	methodFileStop = http.MethodPost
	uriFileStop    = "/player/file/stop"

	methodPlayerState        = http.MethodGet
	uriPlayerState           = "/player/state"
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

	methodRecorderState = http.MethodGet
	uriRecorderState    = "/recorder/state"
	methodRecorderStart = http.MethodPost
	uriRecorderStart    = "/recoder/start"
	methodRecoderStop   = http.MethodPost
	uriRecorderStop     = "/recoder/stop"
)

// NewClient return http client
func NewClient(serverAddr string) Client {
	if !strings.HasPrefix(serverAddr, "http") {
		serverAddr = protocol + "://" + serverAddr
	}
	return &client{
		cli:               &fasthttp.Client{},
		filePlayTransport: NewFilePlayTransport(methodFilePlay, serverAddr+uriFilePlay),
		fileStopTransport: NewFileStopTransport(methodFileStop, serverAddr+uriFileStop),
	}
}
