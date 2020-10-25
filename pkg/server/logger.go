package server

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	server Server
	logger log.Logger
}

func (l *loggerMiddleware) FilePlay(ctx context.Context, file, playerIP, playerPort, playerDeviceName string) (uuid string, channels uint16, rate uint32, bitsPerSample uint16, err error) {
	l.logger.Log("FilePlay", "start")
	if uuid, channels, rate, bitsPerSample, err = l.server.FilePlay(ctx, file, playerIP, playerPort, playerDeviceName); err != nil {
		l.logger.Log(
			"FilePlay", "err",
			"file", file,
			"playerIP", playerIP,
			"playerPort", playerPort,
			"playerDeviceName", playerDeviceName,
			"err", err,
		)
		return
	}
	l.logger.Log(
		"FilePlay", "end",
		"uuid", uuid,
		"channels", channels,
		"rate", rate,
		"bitsPerSample", bitsPerSample,
	)
	return
}

func (l *loggerMiddleware) FileStop(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid string) (err error) {
	l.logger.Log("FileStop", "start")
	if err = l.server.FileStop(ctx, playerIP, playerPort, playerDeviceName, uuid); err != nil {
		l.logger.Log(
			"FileStop", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"playerDeviceName", playerDeviceName,
			"uuid", uuid,
			"err", err,
		)
		return
	}
	l.logger.Log("FileStop", "end")
	return
}

func (l *loggerMiddleware) PlayerState(ctx context.Context, playerIP string) (ports, storages, devices []string, err error) {
	l.logger.Log("PlayerState", "start")
	if ports, storages, devices, err = l.server.PlayerState(ctx, playerIP); err != nil {
		l.logger.Log(
			"PlayerReceiveStart", "err",
			"playerIP", playerIP,
			"err", err,
		)
	}
	l.logger.Log(
		"PlayerReceiveStart", "end",
		"ports", ports,
		"storages", storages,
		"devices", devices,
	)
	return
}

func (l *loggerMiddleware) PlayerReceiveStart(ctx context.Context, playerIP, playerPort string, uuid *string) (sUUID string, err error) {
	l.logger.Log("PlayerReceiveStart", "start")
	if sUUID, err = l.server.PlayerReceiveStart(ctx, playerIP, playerPort, uuid); err != nil {
		l.logger.Log(
			"PlayerReceiveStart", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"uuid", *uuid,
			"err", err,
		)
	}
	l.logger.Log(
		"PlayerReceiveStart", "end",
		"uuid", sUUID,
	)
	return
}

func (l *loggerMiddleware) PlayerReceiveStop(ctx context.Context, playerIP, playerPort string) (err error) {
	l.logger.Log("PlayerReceiveStop", "start")
	if err = l.server.PlayerReceiveStop(ctx, playerIP, playerPort); err != nil {
		l.logger.Log(
			"PlayerReceiveStop", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"err", err,
		)
	}
	l.logger.Log("PlayerReceiveStop", "end")
	return
}

func (l *loggerMiddleware) PlayerPlay(ctx context.Context, playerIP, uuid, playerDeviceName string, channels, rate, bitsPerSample uint32) (err error) {
	l.logger.Log("PlayerPlay", "start")
	if err = l.server.PlayerPlay(ctx, playerIP, uuid, playerDeviceName, channels, rate, bitsPerSample); err != nil {
		l.logger.Log(
			"PlayerPlay", "err",
			"playerIP", playerIP,
			"uuid", uuid,
			"playerDeviceName", playerDeviceName,
			"channels", channels,
			"rate", rate,
			"bitsPerSample", bitsPerSample,
			"err", err,
		)
	}
	l.logger.Log("PlayerPlay", "end")
	return
}

func (l *loggerMiddleware) PlayerStop(ctx context.Context, playerIP, playerDeviceName string) (err error) {
	l.logger.Log("PlayerStop", "start")
	if err = l.server.PlayerStop(ctx, playerIP, playerDeviceName); err != nil {
		l.logger.Log(
			"PlayerStop", "err",
			"playerIP", playerIP,
			"playerDeviceName", playerDeviceName,
			"err", err,
		)
	}
	l.logger.Log("PlayerStop", "end")
	return
}

func (l *loggerMiddleware) PlayerClearStorage(ctx context.Context, playerIP, uuid string) (err error) {
	l.logger.Log("PlayerClearStorage", "start")
	if err = l.server.PlayerClearStorage(ctx, playerIP, uuid); err != nil {
		l.logger.Log(
			"PlayerClearStorage", "err",
			"playerIP", playerIP,
			"uuid", uuid,
			"err", err,
		)
	}
	l.logger.Log("PlayerClearStorage", "end")
	return
}

func (l *loggerMiddleware) StartFileRecoding(ctx context.Context, recorderIP, recorderDeviceName string, channels, rate uint32, receivePort, file string) (err error) {
	l.logger.Log("StartFileRecoding", "start")
	if err = l.server.StartFileRecoding(ctx, recorderIP, recorderDeviceName, channels, rate, receivePort, file); err != nil {
		l.logger.Log(
			"StartFileRecoding", "err",
			"recorderIP", recorderIP,
			"recorderDeviceName", recorderDeviceName,
			"channels", channels,
			"rate", rate,
			"receivePort", receivePort,
			"file", file,
			"err", err,
		)
	}
	l.logger.Log("StartFileRecoding", "end")
	return
}

func (l *loggerMiddleware) StopFileRecoding(ctx context.Context, recorderIP, recorderDeviceName, receivePort string) (err error) {
	l.logger.Log("StopFileRecoding", "start")
	if err = l.server.StopFileRecoding(ctx, recorderIP, recorderDeviceName, receivePort); err != nil {
		l.logger.Log(
			"StopFileRecoding", "err",
			"recorderIP", recorderIP,
			"recorderDeviceName", recorderDeviceName,
			"receivePort", receivePort,
			"err", err,
		)
	}
	l.logger.Log("StopFileRecoding", "end")
	return
}

func (l *loggerMiddleware) PlayFromRecorder(ctx context.Context, playerIP, playerPort, playerDeviceName string, channels, rate uint32, recorderIP, recorderDeviceName string) (uuid string, err error) {
	l.logger.Log("PlayFromRecorder", "start")
	if uuid, err = l.server.PlayFromRecorder(ctx, playerIP, playerPort, playerDeviceName, channels, rate, recorderIP, recorderDeviceName); err != nil {
		l.logger.Log(
			"PlayFromRecorder", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"playerDeviceName", playerDeviceName,
			"channels", channels,
			"rate", rate,
			"recorderIP", recorderIP,
			"recorderDeviceName", recorderDeviceName,
			"err", err,
		)
	}
	l.logger.Log(
		"PlayFromRecorder", "end",
		"uuid", uuid,
	)
	return
}

func (l *loggerMiddleware) StopFromRecorder(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string) (err error) {
	l.logger.Log("StopFromRecorder", "start")
	if err = l.server.StopFromRecorder(ctx, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName); err != nil {
		l.logger.Log(
			"StopFromRecorder", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"playerDeviceName", playerDeviceName,
			"uuid", uuid,
			"recorderIP", recorderIP,
			"recorderDeviceName", recorderDeviceName,
			"err", err,
		)
	}
	l.logger.Log("StopFromRecorder", "end")
	return
}

func (l *loggerMiddleware) RecorderState(ctx context.Context, recorderIP string) (devices []string, err error) {
	l.logger.Log("PlayerState", "start")
	if devices, err = l.server.RecorderState(ctx, recorderIP); err != nil {
		l.logger.Log(
			"PlayerReceiveStart", "err",
			"recorderIP", recorderIP,
			"err", err,
		)
	}
	l.logger.Log(
		"PlayerReceiveStart", "end",
		"devices", devices,
	)
	return
}

func (l *loggerMiddleware) RecorderStart(ctx context.Context, recorderIP, recorderDeviceName string, channels, rate uint32, dstAddr string) (err error) {
	l.logger.Log("RecorderStart", "start")
	if err = l.server.RecorderStart(ctx, recorderIP, recorderDeviceName, channels, rate, dstAddr); err != nil {
		l.logger.Log(
			"RecorderStart", "err",
			"recorderIP", recorderIP,
			"recorderDeviceName", recorderDeviceName,
			"channels", channels,
			"rate", rate,
			"dstAddr", dstAddr,
			"err", err,
		)
	}
	l.logger.Log("RecorderStart", "end")
	return
}

func (l *loggerMiddleware) RecorderStop(ctx context.Context, recorderIP, recorderDeviceName string) (err error) {
	l.logger.Log("RecorderStop", "start")
	if err = l.server.RecorderStop(ctx, recorderIP, recorderDeviceName); err != nil {
		l.logger.Log(
			"RecorderStop", "err",
			"recorderIP", recorderIP,
			"recorderDeviceName", recorderDeviceName,
			"err", err,
		)
	}
	l.logger.Log("RecorderStop", "end")
	return
}

// NewLoggerMiddleware logger middleware for server.
func NewLoggerMiddleware(server Server, logger log.Logger) Server {
	return &loggerMiddleware{
		server: server,
		logger: logger,
	}
}
