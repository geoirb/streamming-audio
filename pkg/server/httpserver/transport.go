package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

// FilePlayingTransport ...
type FilePlayingTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (file, playerIP, playerPort, playerDeviceName string, err error)
	Encode(res *fasthttp.Response, uuid string, channels uint16, rate uint32) (err error)
}

type filePlayingTransport struct{}

type filePlayingRequest struct {
	File             string `json:"file"`
	PlayerIP         string `json:"playerIP"`
	PlayerPort       string `json:"playerPort"`
	PlayerDeviceName string `json:"playerDeviceName"`
}

func (t *filePlayingTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, string, string, error) {
	var request filePlayingRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.File, request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, err
}

type filePlayingResponse struct {
	UUID     string `json:"uuid"`
	Channels uint16 `json:"channels"`
	Rate     uint32 `json:"rate"`
}

func (t *filePlayingTransport) Encode(res *fasthttp.Response, uuid string, channels uint16, rate uint32) (err error) {
	response := &filePlayingResponse{
		UUID:     uuid,
		Channels: channels,
		Rate:     rate,
	}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newFilePlayingTransport() FilePlayingTransport {
	return &filePlayingTransport{}
}

// PlayerReceiveStartTransport ...
type PlayerReceiveStartTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (playerIP, playerPort string, uuid *string, err error)
	Encode(res *fasthttp.Response, uuid string) (err error)
}

type playerReceiveStartTransport struct{}

type playerReceiveStartRequest struct {
	PlayerIP   string  `json:"playerIP"`
	PlayerPort string  `json:"playerPort"`
	UUID       *string `json:"uuid"`
}

func (t *playerReceiveStartTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, *string, error) {
	var request playerReceiveStartRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, request.UUID, err
}

type playerReceiveStartResponse struct {
	UUID string `json:"uuid"`
}

func (t *playerReceiveStartTransport) Encode(res *fasthttp.Response, uuid string) (err error) {
	response := &playerReceiveStartResponse{
		UUID: uuid,
	}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newPlayerReceiveStartTransport() PlayerReceiveStartTransport {
	return &playerReceiveStartTransport{}
}
