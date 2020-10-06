package server

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	server Server
	logger log.Logger
}

func (l *loggerMiddleware) FilePlaying(ctx context.Context, file, playerIP, playerPort, playerDeviceName string) (uuid string, channels uint16, rate uint32, err error) {
	l.logger.Log("FilePlaying", "start")
	if uuid, channels, rate, err = l.server.FilePlaying(ctx, file, playerIP, playerPort, playerDeviceName); err != nil {
		l.logger.Log(
			"FilePlaying", "err",
			"file", file,
			"playerIP", playerIP,
			"playerPort", playerPort,
			"playerDeviceName", playerDeviceName,
			"err", err,
		)
		return
	}
	l.logger.Log(
		"FilePlaying", "end",
		"uuid", uuid,
		"channels", channels,
		"rate", rate,
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

func (l *loggerMiddleware) PlayerPlay(ctx context.Context, playerIP, uuid, playerDeviceName string, channels, rate uint32) (err error) {
	l.logger.Log("PlayerPlay", "start")
	if err = l.server.PlayerPlay(ctx, playerIP, uuid, playerDeviceName, channels, rate); err != nil {
		l.logger.Log(
			"PlayerPlay", "err",
			"playerIP", playerIP,
			"uuid", uuid,
			"playerDeviceName", playerDeviceName,
			"channels", channels,
			"rate", rate,
			"err", err,
		)
	}
	l.logger.Log("PlayerPlay", "end")
	return
}

func (l *loggerMiddleware) PlayerPause(ctx context.Context, playerIP, playerDeviceName string) (err error) {
	l.logger.Log("PlayerPause", "start")
	if err = l.server.PlayerPause(ctx, playerIP, playerDeviceName); err != nil {
		l.logger.Log(
			"PlayerPause", "err",
			"playerIP", playerIP,
			"playerDeviceName", playerDeviceName,
			"err", err,
		)
	}
	l.logger.Log("PlayerPause", "end")
	return
}

func (l *loggerMiddleware) PlayerStop(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid string) (err error) {
	l.logger.Log("PlayerStop", "start")
	if err = l.server.PlayerStop(ctx, playerIP, playerPort, playerDeviceName, uuid); err != nil {
		l.logger.Log(
			"PlayerStop", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"playerDeviceName", playerDeviceName,
			"uuid", uuid,
			"err", err,
		)
	}
	l.logger.Log("PlayerStop", "end")
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

func (l *loggerMiddleware) RecoderStop(ctx context.Context, recorderIP, recorderDeviceName string) (err error) {
	l.logger.Log("RecoderStop", "start")
	if err = l.server.RecoderStop(ctx, recorderIP, recorderDeviceName); err != nil {
		l.logger.Log(
			"RecoderStop", "err",
			"recorderIP", recorderIP,
			"recorderDeviceName", recorderDeviceName,
			"err", err,
		)
	}
	l.logger.Log("RecoderStop", "end")
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

// NewLoggerMiddleware logger middleware for server
func NewLoggerMiddleware(logger log.Logger, server Server) Server {
	return &loggerMiddleware{
		server: server,
		logger: logger,
	}
}
