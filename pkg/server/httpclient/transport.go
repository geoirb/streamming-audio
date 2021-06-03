package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"
)

// FilePlayTransport ...
type FilePlayTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, file, playerIP, playerPort, playerDeviceName string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (uuid string, channels uint16, rate uint32, bitsPerSample uint16, err error)
}

type filePlayTransport struct {
	method       string
	pathTemplate string
}

type filePlayRequest struct {
	File             string `json:"file"`
	PlayerIP         string `json:"playerIP"`
	PlayerPort       string `json:"playerPort"`
	PlayerDeviceName string `json:"playerDeviceName"`
}

func (t *filePlayTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, file, playerIP, playerPort, playerDeviceName string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := filePlayRequest{
		File:             file,
		PlayerIP:         playerIP,
		PlayerPort:       playerPort,
		PlayerDeviceName: playerDeviceName,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

type filePlayResponse struct {
	UUID          string `json:"uuid"`
	Channels      uint16 `json:"channels"`
	Rate          uint32 `json:"rate"`
	BitsPerSample uint16 `json:"bitsPerSample"`
}

func (t *filePlayTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (uuid string, channels uint16, rate uint32, bitsPerSample uint16, err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
		return
	}

	var response filePlayResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return
	}

	uuid, channels, rate = response.UUID, response.Channels, response.Rate
	return
}

// NewFilePlayTransport ...
func NewFilePlayTransport(method, pathTemplate string) FilePlayTransport {
	return &filePlayTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// FileStopTransport ...
type FileStopTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort, playerDeviceName, uuid string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type fileStopTransport struct {
	pathTemplate string
	method       string
}

type fileStopRequest struct {
	PlayerIP         string `json:"playerIP"`
	PlayerPort       string `json:"playerPort"`
	PlayerDeviceName string `json:"playerDeviceName"`
	UUID             string `json:"uuid"`
}

func (t *fileStopTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort, playerDeviceName, uuid string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := fileStopRequest{
		PlayerIP:         playerIP,
		PlayerPort:       playerPort,
		PlayerDeviceName: playerDeviceName,
		UUID:             uuid,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *fileStopTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewFileStopTransport ...
func NewFileStopTransport(method, pathTemplate string) FileStopTransport {
	return &fileStopTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// PlayerStateTransport ...
type PlayerStateTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (ports, storages, devices []string, err error)
}

type playerStateTransport struct {
	method       string
	pathTemplate string
}

type playerStateRequest struct {
	PlayerIP string `json:"playerIP"`
}

func (t *playerStateTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := playerStateRequest{
		PlayerIP: playerIP,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

type playerStateResponse struct {
	Ports    []string `json:"ports"`
	Storages []string `json:"storages"`
	Devices  []string `json:"devices"`
}

func (t *playerStateTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (ports, storages, devices []string, err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
		return
	}

	var response playerStateResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return
	}

	ports, storages, devices = response.Ports, response.Storages, response.Devices
	return
}

// NewPlayerStateTransport ...
func NewPlayerStateTransport(method, pathTemplate string) PlayerStateTransport {
	return &playerStateTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// PlayerReceiveStartTransport ...
type PlayerReceiveStartTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort string, uuid *string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (uuid string, err error)
}

type playerReceiveStartTransport struct {
	method       string
	pathTemplate string
}

type playerReceiveStartRequest struct {
	PlayerIP   string  `json:"playerIP"`
	PlayerPort string  `json:"playerPort"`
	UUID       *string `json:"uuid,omitempty"`
}

func (t *playerReceiveStartTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort string, uuid *string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := playerReceiveStartRequest{
		PlayerIP:   playerIP,
		PlayerPort: playerPort,
		UUID:       uuid,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

type playerReceiveStartResponse struct {
	UUID string `json:"uuid"`
}

func (t *playerReceiveStartTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (uuid string, err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
		return
	}

	var response playerReceiveStartResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return
	}

	uuid = response.UUID
	return
}

// NewPlayerReceiveStartTransport ...
func NewPlayerReceiveStartTransport(method, pathTemplate string) PlayerReceiveStartTransport {
	return &playerReceiveStartTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// PlayerReceiveStopTransport ...
type PlayerReceiveStopTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type playerReceiveStopTransport struct {
	method       string
	pathTemplate string
}

type playerReceiveStopRequest struct {
	PlayerIP   string `json:"playerIP"`
	PlayerPort string `json:"playerPort"`
}

func (t *playerReceiveStopTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := playerReceiveStopRequest{
		PlayerIP:   playerIP,
		PlayerPort: playerPort,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *playerReceiveStopTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewPlayerReceiveStopTransport ...
func NewPlayerReceiveStopTransport(method, pathTemplate string) PlayerReceiveStopTransport {
	return &playerReceiveStopTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// PlayerPlayTransport ...
type PlayerPlayTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, uuid, playerDeviceName string, channels, rate, bitsPerSample uint32) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type playerPlayTransport struct {
	method       string
	pathTemplate string
}

type playerPlayRequest struct {
	PlayerIP         string `json:"playerIP"`
	UUID             string `json:"uuid"`
	PlayerDeviceName string `json:"playerDeviceName"`
	Channels         uint32 `json:"channels"`
	Rate             uint32 `json:"rate"`
	BitsPerSample    uint32 `json:"bitsPerSample"`
}

func (t *playerPlayTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, uuid, playerDeviceName string, channels, rate, bitsPerSample uint32) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := playerPlayRequest{
		PlayerIP:         playerIP,
		UUID:             uuid,
		PlayerDeviceName: playerDeviceName,
		Channels:         channels,
		Rate:             rate,
		BitsPerSample:    bitsPerSample,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *playerPlayTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewPlayerPlayTransport ...
func NewPlayerPlayTransport(method, pathTemplate string) PlayerPlayTransport {
	return &playerPlayTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// PlayerStopTransport ...
type PlayerStopTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerDeviceName string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type playerStopTransport struct {
	method       string
	pathTemplate string
}

type playerStopTransportRequest struct {
	PlayerIP         string `json:"playerIP"`
	PlayerDeviceName string `json:"playerDeviceName"`
}

func (t *playerStopTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerDeviceName string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := playerStopTransportRequest{
		PlayerIP:         playerIP,
		PlayerDeviceName: playerDeviceName,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *playerStopTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewPlayerStopTransport ...
func NewPlayerStopTransport(method, pathTemplate string) PlayerStopTransport {
	return &playerStopTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// PlayerClearStorageTransport ...
type PlayerClearStorageTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, uuid string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type playerClearStorageTransport struct {
	method       string
	pathTemplate string
}

type playerClearStorageRequest struct {
	PlayerIP string `json:"playerIP"`
	UUID     string `json:"uuid"`
}

func (t *playerClearStorageTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, uuid string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := playerClearStorageRequest{
		PlayerIP: playerIP,
		UUID:     uuid,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *playerClearStorageTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewPlayerClearStorageTransport ...
func NewPlayerClearStorageTransport(method, pathTemplate string) PlayerClearStorageTransport {
	return &playerClearStorageTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// StartFileRecordingTransport ...
type StartFileRecordingTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP, recorderDeviceName string, channels, rate uint32, receivePort, file string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type startFileRecordingTransport struct {
	method       string
	pathTemplate string
}

type startFileRecordingRequest struct {
	RecorderIP         string `json:"playerIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	Channels           uint32 `json:"channels"`
	Rate               uint32 `json:"rate"`
	ReceivePort        string `json:"receivePort"`
	File               string `json:"file"`
}

func (t *startFileRecordingTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP, recorderDeviceName string, channels, rate uint32, receivePort, file string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := startFileRecordingRequest{
		RecorderIP:         recorderIP,
		RecorderDeviceName: recorderDeviceName,
		Channels:           channels,
		Rate:               rate,
		ReceivePort:        receivePort,
		File:               file,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *startFileRecordingTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewStartFileRecordingTransport ...
func NewStartFileRecordingTransport(method, pathTemplate string) StartFileRecordingTransport {
	return &startFileRecordingTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// StopFileRecordingTransport ...
type StopFileRecordingTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP, recorderDeviceName, receivePort string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type stopFileRecordingTransport struct {
	method       string
	pathTemplate string
}

type stopFileRecordingRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	ReceivePort        string `json:"receivePort"`
}

func (t *stopFileRecordingTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP, recorderDeviceName, receivePort string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := stopFileRecordingRequest{
		RecorderIP:         recorderIP,
		RecorderDeviceName: recorderDeviceName,
		ReceivePort:        receivePort,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *stopFileRecordingTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewStopFileRecordingTransport ...
func NewStopFileRecordingTransport(method, pathTemplate string) StopFileRecordingTransport {
	return &stopFileRecordingTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// PlayFromRecorderTransport ...
type PlayFromRecorderTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort, playerDeviceName string, channels, rate uint32, recorderIP, recorderDeviceName string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (uuid string, err error)
}

type playFromRecorderTransport struct {
	method       string
	pathTemplate string
}

type playFromRecorderRequest struct {
	PlayerIP           string `json:"playerIP"`
	PlayerPort         string `json:"playerPort"`
	PlayerDeviceName   string `json:"playerDeviceName"`
	Channels           uint32 `json:"channels"`
	Rate               uint32 `json:"rate"`
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
}

func (t *playFromRecorderTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort, playerDeviceName string, channels, rate uint32, recorderIP, recorderDeviceName string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := playFromRecorderRequest{
		PlayerIP:           playerIP,
		PlayerPort:         playerPort,
		PlayerDeviceName:   playerDeviceName,
		Channels:           channels,
		Rate:               rate,
		RecorderIP:         recorderIP,
		RecorderDeviceName: recorderDeviceName,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

type playFromRecorderResponse struct {
	UUID string `json:"uuid"`
}

func (t *playFromRecorderTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (uuid string, err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
		return
	}

	var response playFromRecorderResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return
	}

	uuid = response.UUID
	return
}

// NewPlayFromRecorderTransport ...
func NewPlayFromRecorderTransport(method, pathTemplate string) PlayFromRecorderTransport {
	return &playFromRecorderTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// StopFromRecorderTransport ...
type StopFromRecorderTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type stopFromRecorderTransport struct {
	method       string
	pathTemplate string
}

type stopFromRecorderRequest struct {
	PlayerIP           string `json:"playerIP"`
	PlayerPort         string `json:"playerPort"`
	PlayerDeviceName   string `json:"playerDeviceName"`
	UUID               string `json:"uuid"`
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
}

func (t *stopFromRecorderTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := stopFromRecorderRequest{
		PlayerIP:           playerIP,
		PlayerPort:         playerPort,
		PlayerDeviceName:   playerDeviceName,
		UUID:               uuid,
		RecorderIP:         recorderIP,
		RecorderDeviceName: recorderDeviceName,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *stopFromRecorderTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewStopFromRecorderTransport ...
func NewStopFromRecorderTransport(method, pathTemplate string) StopFromRecorderTransport {
	return &stopFromRecorderTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// RecorderStateTransport ...
type RecorderStateTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (devices []string, err error)
}

type recorderStateTransport struct {
	method       string
	pathTemplate string
}

type recorderStateRequest struct {
	RecorderIP string `json:"recorderIP"`
}

func (t *recorderStateTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := recorderStateRequest{
		RecorderIP: recorderIP,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

type recorderStateResponse struct {
	Devices []string `json:"devices"`
}

func (t *recorderStateTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (devices []string, err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
		return
	}

	var response recorderStateResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return
	}

	devices = response.Devices
	return
}

// NewRecorderStateTransport ...
func NewRecorderStateTransport(method, pathTemplate string) RecorderStateTransport {
	return &recorderStateTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// RecorderStartTransport ...
type RecorderStartTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP, recorderDeviceName string, channels, rate uint32, dstAddr string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type recorderStartTransport struct {
	method       string
	pathTemplate string
}

type recorderStartRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
	Channels           uint32 `json:"channels"`
	Rate               uint32 `json:"rate"`
	DstAddr            string `json:"dstAddr"`
}

func (t *recorderStartTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP, recorderDeviceName string, channels, rate uint32, dstAddr string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := recorderStartRequest{
		RecorderIP:         recorderIP,
		RecorderDeviceName: recorderDeviceName,
		Channels:           channels,
		Rate:               rate,
		DstAddr:            dstAddr,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *recorderStartTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewRecorderStartTransport ...
func NewRecorderStartTransport(method, pathTemplate string) RecorderStartTransport {
	return &recorderStartTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}

// RecorderStopTransport ...
type RecorderStopTransport interface {
	EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP, recorderDeviceName string) (err error)
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error)
}

type recorderStopTransport struct {
	method       string
	pathTemplate string
}

type recorderStopRequest struct {
	RecorderIP         string `json:"recorderIP"`
	RecorderDeviceName string `json:"recorderDeviceName"`
}

func (t *recorderStopTransport) EncodeRequest(ctx context.Context, req *fasthttp.Request, recorderIP, recorderDeviceName string) (err error) {
	req.Header.SetMethod(t.method)
	req.SetRequestURI(t.pathTemplate)

	request := recorderStopRequest{
		RecorderIP:         recorderIP,
		RecorderDeviceName: recorderDeviceName,
	}
	body, err := json.Marshal(&request)
	if err != nil {
		return
	}

	req.SetBody(body)
	return
}

func (t *recorderStopTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
	}
	return
}

// NewRecorderStopTransport ...
func NewRecorderStopTransport(method, pathTemplate string) RecorderStopTransport {
	return &recorderStopTransport{
		method:       method,
		pathTemplate: pathTemplate,
	}
}
