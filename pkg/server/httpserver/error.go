package httpserver

import (
	"net/http"

	"github.com/valyala/fasthttp"

	"github.com/geoirb/audio-service/pkg/server"
)

const (
	codeDeviceIsBusy   = http.StatusInternalServerError
	codeDeviceNotFound = http.StatusNotFound
	codePortIsBusy     = http.StatusInternalServerError
	codePortNotFound   = http.StatusNotFound
)

type errorProcessing func(res *fasthttp.Response, err error, statusCode int)

// ErrorProcessing ...
func ErrorProcessing(res *fasthttp.Response, err error, statusCode int) {
	res.SetBody([]byte(err.Error()))
	res.SetStatusCode(statusCode)

	switch err {
	case server.ErrDeviceIsBusy:
		res.SetStatusCode(codeDeviceIsBusy)
	case server.ErrDeviceNotFound:
		res.SetStatusCode(codeDeviceNotFound)
	case server.ErrPortIsBusy:
		res.SetStatusCode(codePortIsBusy)
	case server.ErrPortNotFound:
		res.SetStatusCode(codePortNotFound)
	default:
		res.SetStatusCode(http.StatusInternalServerError)
	}
	if statusCode != -1 {
		res.SetStatusCode(statusCode)
	}
}
