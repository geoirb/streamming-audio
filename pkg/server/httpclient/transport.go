// todo
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
	DecodeResponse(ctx context.Context, res *fasthttp.Response) (uuid string, channels uint16, rate uint32, err error)
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
	UUID     string `json:"uuid"`
	Channels uint16 `json:"channels"`
	Rate     uint32 `json:"rate"`
}

func (t *filePlayTransport) DecodeResponse(ctx context.Context, res *fasthttp.Response) (uuid string, channels uint16, rate uint32, err error) {
	if res.StatusCode() != http.StatusOK {
		err = fmt.Errorf(string(res.Body()))
		return
	}

	var response filePlayResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return
	}

	uuid = response.UUID
	channels = response.Channels
	rate = response.Rate
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
		return
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
