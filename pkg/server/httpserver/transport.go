package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

// FilePlayTransport ...
type FilePlayTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (file, playerIP, playerPort, playerDeviceName string, err error)
	Encode(res *fasthttp.Response, uuid string, channels uint16, rate uint32) (err error)
}

type filePlayTransport struct{}

type filePlayRequest struct {
	File             string `json:"file"`
	PlayerIP         string `json:"playerIP"`
	PlayerPort       string `json:"playerPort"`
	PlayerDeviceName string `json:"playerDeviceName"`
}

func (t *filePlayTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, string, string, error) {
	var request filePlayRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.File, request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, err
}

type filePlayResponse struct {
	UUID     string `json:"uuid"`
	Channels uint16 `json:"channels"`
	Rate     uint32 `json:"rate"`
}

func (t *filePlayTransport) Encode(res *fasthttp.Response, uuid string, channels uint16, rate uint32) (err error) {
	response := &filePlayResponse{
		UUID:     uuid,
		Channels: channels,
		Rate:     rate,
	}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newFilePlayTransport() FilePlayTransport {
	return &filePlayTransport{}
}

// FileStopTransport ...
type FileStopTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (playerIP, playerPort, playerDeviceName, uuid string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type fileStopTransport struct{}

type fileStopRequest struct {
	PlayerIP         string `json:"playerIP"`
	PlayerPort       string `json:"playerPort"`
	PlayerDeviceName string `json:"playerDeviceName"`
	UUID             string `json:"uuid"`
}

func (t *fileStopTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, string, string, error) {
	var request fileStopRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, request.UUID, err
}

type fileStopResponse struct{}

func (t *fileStopTransport) Encode(res *fasthttp.Response) (err error) {
	response := &fileStopResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newFileStopTransport() FileStopTransport {
	return &fileStopTransport{}
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

// PlayerReceiveStopTransport ...
type PlayerReceiveStopTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (playerIP, playerPort string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type playerReceiveStopTransport struct{}

type playerReceiveStopRequest struct {
	PlayerIP   string `json:"playerIP"`
	PlayerPort string `json:"playerPort"`
}

func (t *playerReceiveStopTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, error) {
	var request playerReceiveStopRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, err
}

type playerReceiveStopResponse struct{}

func (t *playerReceiveStopTransport) Encode(res *fasthttp.Response) (err error) {
	response := &playerReceiveStopResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newPlayerReceiveStopTransport() PlayerReceiveStopTransport {
	return &playerReceiveStopTransport{}
}

// PlayerPlayTransport ...
type PlayerPlayTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (playerIP, uuid, playerDeviceName string, channels, rate uint32, err error)
	Encode(res *fasthttp.Response) (err error)
}

type playerPlayTransport struct{}

type playerPlayRequest struct {
	PlayerIP         string `json:"playerIP"`
	UUID             string `json:"uuid"`
	PlayerDeviceName string `json:"playerDeviceName"`
	Channels         uint32 `json:"channels"`
	Rate             uint32 `json:"rate"`
}

func (t *playerPlayTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, string, uint32, uint32, error) {
	var request playerPlayRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.UUID, request.PlayerDeviceName, request.Channels, request.Rate, err
}

type playerPlayResponse struct{}

func (t *playerPlayTransport) Encode(res *fasthttp.Response) (err error) {
	response := &playerPlayResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newPlayerPlayTransport() PlayerPlayTransport {
	return &playerPlayTransport{}
}

// PlayerStopTransport ...
type PlayerStopTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (playerIP, playerDeviceName string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type playerStopTransport struct{}

type playerStopRequest struct {
	PlayerIP         string `json:"playerIP"`
	PlayerDeviceName string `json:"playerDeviceName"`
}

func (t *playerStopTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, error) {
	var request playerStopRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerDeviceName, err
}

type playerStopResponse struct{}

func (t *playerStopTransport) Encode(res *fasthttp.Response) (err error) {
	response := &playerStopResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newPlayerStopTransport() PlayerStopTransport {
	return &playerStopTransport{}
}

// PlayerClearStorageTransport ...
type PlayerClearStorageTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (playerIP, uuid string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type playerClearStorageTransport struct{}

type playerClearStorageRequest struct {
	PlayerIP string `json:"playerIP"`
	UUID     string `json:"uuid"`
}

func (t *playerClearStorageTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, error) {
	var request playerClearStorageRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.UUID, err
}

type playerClearStorageResponse struct{}

func (t *playerClearStorageTransport) Encode(res *fasthttp.Response) (err error) {
	response := &playerClearStorageResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newPlayerClearStorageTransport() PlayerClearStorageTransport {
	return &playerClearStorageTransport{}
}

// StartFileRecodingTransport ...
type StartFileRecodingTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (recorderIP, recorderDeviceName string, channels, rate uint32, receivePort, file string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type startFileRecodingTransport struct{}

type startFileRecodingRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	Channels           uint32 `json:"channels"`
	Rate               uint32 `json:"rate"`
	ReceivePort        string `json:"receivePort"`
	File               string `json:"file"`
}

func (t *startFileRecodingTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, uint32, uint32, string, string, error) {
	var request startFileRecodingRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, request.RecorderDeviceName, request.Channels, request.Rate, request.ReceivePort, request.File, err
}

type startFileRecodingResponse struct{}

func (t *startFileRecodingTransport) Encode(res *fasthttp.Response) (err error) {
	response := &startFileRecodingResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newStartFileRecodingTransport() StartFileRecodingTransport {
	return &startFileRecodingTransport{}
}

// StopFileRecodingTransport ...
type StopFileRecodingTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (recorderIP, recorderDeviceName, receivePort string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type stopFileRecodingTransport struct{}

type stopFileRecodingRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	ReceivePort        string `json:"receivePort"`
}

func (t *stopFileRecodingTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, string, error) {
	var request stopFileRecodingRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, request.RecorderDeviceName, request.ReceivePort, err
}

type stopFileRecodingResponse struct{}

func (t *stopFileRecodingTransport) Encode(res *fasthttp.Response) (err error) {
	response := &stopFileRecodingResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newStopFileRecodingTransport() StopFileRecodingTransport {
	return &stopFileRecodingTransport{}
}

// PlayFromRecorderTransport ...
type PlayFromRecorderTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (playerIP, playerPort, playerDeviceName string, channels, rate uint32, recorderIP, recorderDeviceName string, err error)
	Encode(res *fasthttp.Response, uuid string) (err error)
}

type playFromRecorderTransport struct{}

type playFromRecorderRequest struct {
	PlayerIP           string `json:"playerIP"`
	PlayerPort         string `json:"playerPort"`
	PlayerDeviceName   string `json:"playerDeviceName"`
	Channels           uint32 `json:"channels"`
	Rate               uint32 `json:"rate"`
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
}

func (t *playFromRecorderTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, string, uint32, uint32, string, string, error) {
	var request playFromRecorderRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, request.Channels, request.Rate, request.RecorderIP, request.RecorderDeviceName, err
}

type playFromRecorderResponse struct {
	UUID string `json:"uuid"`
}

func (t *playFromRecorderTransport) Encode(res *fasthttp.Response, uuid string) (err error) {
	response := &playFromRecorderResponse{
		UUID: uuid,
	}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newPlayFromRecorderTransport() PlayFromRecorderTransport {
	return &playFromRecorderTransport{}
}

// StopFromRecorderTransport ...
type StopFromRecorderTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type stopFromRecorderTransport struct{}

type stopFromRecorderRequest struct {
	PlayerIP           string `json:"playerIP"`
	PlayerPort         string `json:"playerPort"`
	PlayerDeviceName   string `json:"playerDeviceName"`
	UUID               string `json:"uuid"`
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
}

func (t *stopFromRecorderTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, string, string, string, string, error) {
	var request stopFromRecorderRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, request.UUID, request.RecorderIP, request.RecorderDeviceName, err
}

type stopFromRecorderResponse struct{}

func (t *stopFromRecorderTransport) Encode(res *fasthttp.Response) (err error) {
	response := &stopFromRecorderResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newStopFromRecorderTransport() StopFromRecorderTransport {
	return &stopFromRecorderTransport{}
}

// RecorderStartTransport ...
type RecorderStartTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (recorderIP, recorderDeviceName string, channels, rate uint32, dstAddr string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type recorderStartTransport struct{}

type recorderStartRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	Channels           uint32 `json:"channels"`
	Rate               uint32 `json:"rate"`
	DstAddr            string `json:"dstAddr"`
}

func (t *recorderStartTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, uint32, uint32, string, error) {
	var request recorderStartRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, request.RecorderDeviceName, request.Channels, request.Rate, request.DstAddr, err
}

type recorderStartResponse struct{}

func (t *recorderStartTransport) Encode(res *fasthttp.Response) (err error) {
	response := &recorderStartResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newRecorderStartTransport() RecorderStartTransport {
	return &recorderStartTransport{}
}

// RecorderStopTransport ...
type RecorderStopTransport interface {
	Decode(ctx *fasthttp.RequestCtx) (recorderIP, recorderDeviceName string, err error)
	Encode(res *fasthttp.Response) (err error)
}

type recorderStopTransport struct{}

type recorderStopRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
}

func (t *recorderStopTransport) Decode(ctx *fasthttp.RequestCtx) (string, string, error) {
	var request recorderStopRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, request.RecorderDeviceName, err
}

type recorderStopResponse struct{}

func (t *recorderStopTransport) Encode(res *fasthttp.Response) (err error) {
	response := &recorderStopResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newRecorderStopTransport() RecorderStopTransport {
	return &recorderStopTransport{}
}
