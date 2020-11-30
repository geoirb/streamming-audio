package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

// FilePlayTransport ...
type FilePlayTransport interface {
	DecodeRequest(ctx *fasthttp.RequestCtx) (file, playerIP, playerPort, playerDeviceName string, err error)
	EncodeResponse(res *fasthttp.Response, uuid string, channels uint16, rate uint32, bitsPerSample uint16) (err error)
}

type filePlayTransport struct{}

type filePlayRequest struct {
	File             string `json:"file"`
	PlayerIP         string `json:"playerIP"`
	PlayerPort       string `json:"playerPort"`
	PlayerDeviceName string `json:"playerDeviceName"`
}

func (t *filePlayTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, string, string, error) {
	var request filePlayRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.File, request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, err
}

type filePlayResponse struct {
	UUID          string `json:"uuid"`
	Channels      uint16 `json:"channels"`
	Rate          uint32 `json:"rate"`
	BitsPerSample uint16 `json:"bitsPerSample "`
}

func (t *filePlayTransport) EncodeResponse(res *fasthttp.Response, uuid string, channels uint16, rate uint32, bitsPerSample uint16) (err error) {
	response := &filePlayResponse{
		UUID:          uuid,
		Channels:      channels,
		Rate:          rate,
		BitsPerSample: bitsPerSample,
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
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP, playerPort, playerDeviceName, uuid string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type fileStopTransport struct{}

type fileStopRequest struct {
	PlayerIP         string `json:"playerIP"`
	PlayerPort       string `json:"playerPort"`
	PlayerDeviceName string `json:"playerDeviceName"`
	UUID             string `json:"uuid"`
}

func (t *fileStopTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, string, string, error) {
	var request fileStopRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, request.UUID, err
}

type fileStopResponse struct{}

func (t *fileStopTransport) EncodeResponse(res *fasthttp.Response) (err error) {
	response := &fileStopResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newFileStopTransport() FileStopTransport {
	return &fileStopTransport{}
}

// PlayerStateTransport ...
type PlayerStateTransport interface {
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP string, err error)
	EncodeResponse(res *fasthttp.Response, ports, storages, devices []string) (err error)
}

type playerStateTransport struct{}

type playerStateRequest struct {
	PlayerIP string `json:"playerIP"`
}

func (t *playerStateTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, error) {
	var request playerStateRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, err
}

type playerStateResponse struct {
	Ports    []string `json:"ports"`
	Storages []string `json:"storages"`
	Devices  []string `json:"devices"`
}

func (t *playerStateTransport) EncodeResponse(res *fasthttp.Response, ports, storages, devices []string) (err error) {
	response := &playerStateResponse{
		Ports:    ports,
		Storages: storages,
		Devices:  devices,
	}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newPlayerStateTransport() PlayerStateTransport {
	return &playerStateTransport{}
}

// PlayerReceiveStartTransport ...
type PlayerReceiveStartTransport interface {
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP, playerPort string, uuid *string, err error)
	EncodeResponse(res *fasthttp.Response, uuid string) (err error)
}

type playerReceiveStartTransport struct{}

type playerReceiveStartRequest struct {
	PlayerIP   string  `json:"playerIP"`
	PlayerPort string  `json:"playerPort"`
	UUID       *string `json:"uuid"`
}

func (t *playerReceiveStartTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, *string, error) {
	var request playerReceiveStartRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, request.UUID, err
}

type playerReceiveStartResponse struct {
	UUID string `json:"uuid"`
}

func (t *playerReceiveStartTransport) EncodeResponse(res *fasthttp.Response, uuid string) (err error) {
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
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP, playerPort string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type playerReceiveStopTransport struct{}

type playerReceiveStopRequest struct {
	PlayerIP   string `json:"playerIP"`
	PlayerPort string `json:"playerPort"`
}

func (t *playerReceiveStopTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, error) {
	var request playerReceiveStopRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, err
}

type playerReceiveStopResponse struct{}

func (t *playerReceiveStopTransport) EncodeResponse(res *fasthttp.Response) (err error) {
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
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP, uuid, playerDeviceName string, channels, rate, bitsPerSample uint32, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type playerPlayTransport struct{}

type playerPlayRequest struct {
	PlayerIP         string `json:"playerIP"`
	UUID             string `json:"uuid"`
	PlayerDeviceName string `json:"playerDeviceName"`
	Channels         uint32 `json:"channels"`
	Rate             uint32 `json:"rate"`
	BitsPerSample    uint32 `json:"bitsPerSample"`
}

func (t *playerPlayTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, string, uint32, uint32, uint32, error) {
	var request playerPlayRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.UUID, request.PlayerDeviceName, request.Channels, request.Rate, request.BitsPerSample, err
}

type playerPlayResponse struct{}

func (t *playerPlayTransport) EncodeResponse(res *fasthttp.Response) (err error) {
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
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP, playerDeviceName string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type playerStopTransport struct{}

type playerStopRequest struct {
	PlayerIP         string `json:"playerIP"`
	PlayerDeviceName string `json:"playerDeviceName"`
}

func (t *playerStopTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, error) {
	var request playerStopRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerDeviceName, err
}

type playerStopResponse struct{}

func (t *playerStopTransport) EncodeResponse(res *fasthttp.Response) (err error) {
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
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP, uuid string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type playerClearStorageTransport struct{}

type playerClearStorageRequest struct {
	PlayerIP string `json:"playerIP"`
	UUID     string `json:"uuid"`
}

func (t *playerClearStorageTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, error) {
	var request playerClearStorageRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.UUID, err
}

type playerClearStorageResponse struct{}

func (t *playerClearStorageTransport) EncodeResponse(res *fasthttp.Response) (err error) {
	response := &playerClearStorageResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newPlayerClearStorageTransport() PlayerClearStorageTransport {
	return &playerClearStorageTransport{}
}

// StartFileRecordingTransport ...
type StartFileRecordingTransport interface {
	DecodeRequest(ctx *fasthttp.RequestCtx) (recorderIP, recorderDeviceName string, channels, rate uint32, receivePort, file string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type startFileRecordingTransport struct{}

type startFileRecordingRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	Channels           uint32 `json:"channels"`
	Rate               uint32 `json:"rate"`
	ReceivePort        string `json:"receivePort"`
	File               string `json:"file"`
}

func (t *startFileRecordingTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, uint32, uint32, string, string, error) {
	var request startFileRecordingRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, request.RecorderDeviceName, request.Channels, request.Rate, request.ReceivePort, request.File, err
}

type startFileRecordingResponse struct{}

func (t *startFileRecordingTransport) EncodeResponse(res *fasthttp.Response) (err error) {
	response := &startFileRecordingResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newStartFileRecordingTransport() StartFileRecordingTransport {
	return &startFileRecordingTransport{}
}

// StopFileRecordingTransport ...
type StopFileRecordingTransport interface {
	DecodeRequest(ctx *fasthttp.RequestCtx) (recorderIP, recorderDeviceName, receivePort string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type stopFileRecordingTransport struct{}

type stopFileRecordingRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	ReceivePort        string `json:"receivePort"`
}

func (t *stopFileRecordingTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, string, error) {
	var request stopFileRecordingRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, request.RecorderDeviceName, request.ReceivePort, err
}

type stopFileRecordingResponse struct{}

func (t *stopFileRecordingTransport) EncodeResponse(res *fasthttp.Response) (err error) {
	response := &stopFileRecordingResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newStopFileRecordingTransport() StopFileRecordingTransport {
	return &stopFileRecordingTransport{}
}

// PlayFromRecorderTransport ...
type PlayFromRecorderTransport interface {
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP, playerPort, playerDeviceName string, channels, rate uint32, recorderIP, recorderDeviceName string, err error)
	EncodeResponse(res *fasthttp.Response, uuid string) (err error)
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

func (t *playFromRecorderTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, string, uint32, uint32, string, string, error) {
	var request playFromRecorderRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, request.Channels, request.Rate, request.RecorderIP, request.RecorderDeviceName, err
}

type playFromRecorderResponse struct {
	UUID string `json:"uuid"`
}

func (t *playFromRecorderTransport) EncodeResponse(res *fasthttp.Response, uuid string) (err error) {
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
	DecodeRequest(ctx *fasthttp.RequestCtx) (playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
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

func (t *stopFromRecorderTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, string, string, string, string, error) {
	var request stopFromRecorderRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.PlayerIP, request.PlayerPort, request.PlayerDeviceName, request.UUID, request.RecorderIP, request.RecorderDeviceName, err
}

type stopFromRecorderResponse struct{}

func (t *stopFromRecorderTransport) EncodeResponse(res *fasthttp.Response) (err error) {
	response := &stopFromRecorderResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newStopFromRecorderTransport() StopFromRecorderTransport {
	return &stopFromRecorderTransport{}
}

// RecorderStateTransport ...
type RecorderStateTransport interface {
	DecodeRequest(ctx *fasthttp.RequestCtx) (recorderIP string, err error)
	EncodeResponse(res *fasthttp.Response, devices []string) (err error)
}

type recorderStateTransport struct{}

type recorderStateRequest struct {
	RecorderIP string `json:"recorderIP"`
}

func (t *recorderStateTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, error) {
	var request recorderStateRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, err
}

type recorderStateResponse struct {
	Devices []string `json:"devices"`
}

func (t *recorderStateTransport) EncodeResponse(res *fasthttp.Response, devices []string) (err error) {
	response := &recorderStateResponse{
		Devices: devices,
	}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newRecorderStateTransport() RecorderStateTransport {
	return &recorderStateTransport{}
}

// RecorderStartTransport ...
type RecorderStartTransport interface {
	DecodeRequest(ctx *fasthttp.RequestCtx) (recorderIP, recorderDeviceName string, channels, rate uint32, dstAddr string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type recorderStartTransport struct{}

type recorderStartRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	Channels           uint32 `json:"channels"`
	Rate               uint32 `json:"rate"`
	DstAddr            string `json:"dstAddr"`
}

func (t *recorderStartTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, uint32, uint32, string, error) {
	var request recorderStartRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, request.RecorderDeviceName, request.Channels, request.Rate, request.DstAddr, err
}

type recorderStartResponse struct{}

func (t *recorderStartTransport) EncodeResponse(res *fasthttp.Response) (err error) {
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
	DecodeRequest(ctx *fasthttp.RequestCtx) (recorderIP, recorderDeviceName string, err error)
	EncodeResponse(res *fasthttp.Response) (err error)
}

type recorderStopTransport struct{}

type recorderStopRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
}

func (t *recorderStopTransport) DecodeRequest(ctx *fasthttp.RequestCtx) (string, string, error) {
	var request recorderStopRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	return request.RecorderIP, request.RecorderDeviceName, err
}

type recorderStopResponse struct{}

func (t *recorderStopTransport) EncodeResponse(res *fasthttp.Response) (err error) {
	response := &recorderStopResponse{}
	body, err := json.Marshal(response)
	res.SetBody(body)
	res.SetStatusCode(http.StatusOK)
	return
}

func newRecorderStopTransport() RecorderStopTransport {
	return &recorderStopTransport{}
}
