package httpclient

import (
	"context"

	"github.com/valyala/fasthttp"

	"github.com/geoirb/audio-service/pkg/server"
)

// Client to audio-service
type Client interface {
	server.Server
}

type client struct {
	cli *fasthttp.Client

	filePlayTransport           FilePlayTransport
	fileStopTransport           FileStopTransport
	playerStateTransport        PlayerStateTransport
	playerReceiveStartTransport PlayerReceiveStartTransport
	playerReceiveStopTransport  PlayerReceiveStopTransport
	playerPlayTransport         PlayerPlayTransport
	playerStopTransport         PlayerStopTransport
	playerClearStorageTransport PlayerClearStorageTransport
	startFileRecodingTransport  StartFileRecodingTransport
	stopFileRecodingTransport   StopFileRecodingTransport
	playFromRecorderTransport   PlayFromRecorderTransport
	stopFromRecorderTransport   StopFromRecorderTransport
	recorderStateTransport      RecorderStateTransport
	recorderStartTransport      RecorderStartTransport
	recorderStopTransport       RecorderStopTransport
}

// FilePlay send file to player with playerIP on port and play on playerDeviceName
// channel and rate audio info from file.
// Player save audio from server in storage with uuid.
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

// FileStop stop send file to player with playerIP on port.
// Stop play audio on playerDeviceName on player with playerIP
// Clear storage with uuid on player with playerIP
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

// PlayerState return all busy ports, devices on player and existing storage
func (c *client) PlayerState(ctx context.Context, playerIP string) (ports, storages, devices []string, err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.playerStateTransport.EncodeRequest(ctx, req, playerIP); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.playerStateTransport.DecodeResponse(ctx, res)
}

// PlayerReceiveStart player with playerIP start receive signal from server on playerPort.
// uuid of the storage existing on the player
// if the storage with uuid does not exist or the uuid is nil, a new storage will be created on the player
// The signal will be stored in the storage sUUID
func (c *client) PlayerReceiveStart(ctx context.Context, playerIP, playerPort string, uuid *string) (sUUID string, err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.playerReceiveStartTransport.EncodeRequest(ctx, req, playerIP, playerPort, uuid); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.playerReceiveStartTransport.DecodeResponse(ctx, res)
}

// PlayerReceiveStop player with playerIP stop receive signal from server on playerPort.
func (c *client) PlayerReceiveStop(ctx context.Context, playerIP, playerPort string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.playerReceiveStopTransport.EncodeRequest(ctx, req, playerIP, playerPort); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.playerReceiveStopTransport.DecodeResponse(ctx, res)
}

// PlayerPlay play audio from storage with uuid on player with playerIP on playerDeviceName
// channels, rate - params audio
func (c *client) PlayerPlay(ctx context.Context, playerIP, uuid, playerDeviceName string, channels, rate uint32) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.playerPlayTransport.EncodeRequest(ctx, req, playerIP, uuid, playerDeviceName, channels, rate); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.playerPlayTransport.DecodeResponse(ctx, res)
}

// PlayerStop pause audio on player with playerIP on playerDeviceName
func (c *client) PlayerStop(ctx context.Context, playerIP, playerDeviceName string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.playerStopTransport.EncodeRequest(ctx, req, playerIP, playerDeviceName); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.playerStopTransport.DecodeResponse(ctx, res)
}

// PlayerClearStorage clear storage with uuid on player with playerIP
func (c *client) PlayerClearStorage(ctx context.Context, playerIP, uuid string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.playerClearStorageTransport.EncodeRequest(ctx, req, playerIP, uuid); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.playerClearStorageTransport.DecodeResponse(ctx, res)
}

// StartFileRecoding start receive on receivePort audio signal from recorder with recorderIP from recordeDeviceName and write in file
// channels, rate - params audio
func (c *client) StartFileRecoding(ctx context.Context, recorderIP, recorderDeviceName string, channels, rate uint32, receivePort, file string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.startFileRecodingTransport.EncodeRequest(ctx, req, recorderIP, recorderDeviceName, channels, rate, receivePort, file); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.startFileRecodingTransport.DecodeResponse(ctx, res)
}

// StopFileRecoding stop receive on receivePort audio signal from recorder with recorderIP from recordeDeviceName
func (c *client) StopFileRecoding(ctx context.Context, recorderIP, recorderDeviceName, receivePort string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.stopFileRecodingTransport.EncodeRequest(ctx, req, recorderIP, recorderDeviceName, receivePort); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.stopFileRecodingTransport.DecodeResponse(ctx, res)
}

// PlayFromRecorder play audio on player with playerIP from recorder with recorderIP
func (c *client) PlayFromRecorder(ctx context.Context, playerIP, playerPort, playerDeviceName string, channels, rate uint32, recorderIP, recorderDeviceName string) (uuid string, err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.playFromRecorderTransport.EncodeRequest(ctx, req, playerIP, playerPort, playerDeviceName, channels, rate, recorderIP, recorderDeviceName); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.playFromRecorderTransport.DecodeResponse(ctx, res)
}

// StopFromRecorder stop audio on player with playerIP from recorder with recorderIP
func (c *client) StopFromRecorder(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.stopFromRecorderTransport.EncodeRequest(ctx, req, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.stopFromRecorderTransport.DecodeResponse(ctx, res)
}

// RecorderState return all busy devices on recorder
func (c *client) RecorderState(ctx context.Context, recorderIP string) (devices []string, err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.recorderStateTransport.EncodeRequest(ctx, req, recorderIP); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.recorderStateTransport.DecodeResponse(ctx, res)
}

// RecorderStart start recording audio on recorder with recorderIP from recorderDeviceName and receive on dstAddr
// channels, rate - recoding param
func (c *client) RecorderStart(ctx context.Context, recorderIP, recorderDeviceName string, channels, rate uint32, dstAddr string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.recorderStartTransport.EncodeRequest(ctx, req, recorderIP, recorderDeviceName, channels, rate, dstAddr); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.recorderStartTransport.DecodeResponse(ctx, res)
}

// RecorderStop stop recording audio on recorder with recorderIP from recorderDeviceName
func (c *client) RecorderStop(ctx context.Context, recorderIP, recorderDeviceName string) (err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	if err = c.recorderStopTransport.EncodeRequest(ctx, req, recorderIP, recorderDeviceName); err != nil {
		return
	}

	if err = c.cli.Do(req, res); err != nil {
		return
	}

	return c.recorderStopTransport.DecodeResponse(ctx, res)
}
