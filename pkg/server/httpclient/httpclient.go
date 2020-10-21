// todo
package httpclient

import (
	"context"

	"github.com/valyala/fasthttp"
)

// Client to audio-service
type Client interface {
	// todo
	// server.Server
	FilePlay(ctx context.Context, file, playerIP, playerPort, playerDeviceName string) (uuid string, channels uint16, rate uint32, err error)
	FileStop(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid string) (err error)
}

type client struct {
	cli               *fasthttp.Client
	filePlayTransport FilePlayTransport
	fileStopTransport FileStopTransport
}

func (c *client) FilePlay(ctx context.Context, file, playerIP, playerPort, playerDeviceName string) (uuid string, channels uint16, rate uint32, err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.filePlayTransport.EncodeRequest(ctx, req, file, playerIP, playerPort, playerDeviceName); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.filePlayTransport.DecodeResponse(ctx, res)
}

func (c *client) FileStop(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.fileStopTransport.EncodeRequest(ctx, req, playerIP, playerPort, playerDeviceName, uuid); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.fileStopTransport.DecodeResponse(ctx, res)
}
