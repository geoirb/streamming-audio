package server

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	svc    Server
	logger log.Logger
}

func (l *loggerMiddleware) StartSendingFile(ctx context.Context, destIP, destPort, fileName string) (storageUUID string, channels uint16, rate uint32, err error) {
	l.logger.Log(
		"StartSendingFile", "start",
		"destIP", destIP,
		"destIP", destPort,
		"fileName", fileName,
	)
	if storageUUID, channels, rate, err = l.svc.StartSendingFile(ctx, destIP, destPort, fileName); err != nil {
		l.logger.Log(
			"StartSendingFile", "err",
			"destIP", destIP,
			"destPort", destPort,
			"fileName", fileName,
			"err", err,
		)
	}
	l.logger.Log("StartSendingFile", "end")
	return
}

func (l *loggerMiddleware) StopSending(ctx context.Context, destIP, destPort string) (err error) {
	l.logger.Log(
		"StopSending", "start",
		"destIP", destIP,
		"destIP", destPort,
	)
	if err = l.svc.StopSending(ctx, destIP, destPort); err != nil {
		l.logger.Log(
			"StopSending", "err",
			"destIP", destIP,
			"destPort", destPort,
			"err", err,
		)
	}
	l.logger.Log("StopSending", "end")
	return
}

func (l *loggerMiddleware) StartPlaying(ctx context.Context, playerIP, storageUUID, deviceName string, channels uint16, rate uint32) (err error) {
	l.logger.Log(
		"StartPlaying", "start",
		"playerIP", playerIP,
		"storageUUID", storageUUID,
		"deviceName", deviceName,
		"channels", channels,
		"rate", rate,
	)
	if err = l.svc.StartPlaying(ctx, playerIP, deviceName, storageUUID, channels, rate); err != nil {
		l.logger.Log(
			"StartPlaying", "err",
			"playerIP", playerIP,
			"storageUUID", storageUUID,
			"deviceName", deviceName,
			"channels", channels,
			"rate", rate,
			"err", err,
		)
	}
	l.logger.Log("StartPlaying", "end")
	return
}

func (l *loggerMiddleware) StopPlaying(ctx context.Context, playerIP, deviceName string) (err error) {
	l.logger.Log(
		"StopPlaying", "start",
		"playerIP", playerIP,
		"deviceName", deviceName,
	)
	if err = l.svc.StopPlaying(ctx, playerIP, deviceName); err != nil {
		l.logger.Log(
			"StopPlaying", "err",
			"playerIP", playerIP,
			"deviceName", deviceName,
			"err", err,
		)
	}
	l.logger.Log("StopPlaying", "end")
	return
}

func (l *loggerMiddleware) StartRecordingInFile(ctx context.Context, fileName, receivePort, recoderIP, deviceName string, channels, rate int) (err error) {
	l.logger.Log(
		"StartRecordingInFile", "start",
		"fileName", fileName,
		"receivePort", receivePort,
		"recoderIP", recoderIP,
		"deviceName", deviceName,
		"channels", channels,
		"rate", rate,
	)
	if err = l.svc.StartRecordingInFile(ctx, fileName, receivePort, recoderIP, deviceName, channels, rate); err != nil {
		l.logger.Log(
			"StartRecordingInFile", "err",
			"fileName", fileName,
			"receivePort", receivePort,
			"recoderIP", recoderIP,
			"deviceName", deviceName,
			"channels", channels,
			"rate", rate,
			"err", err,
		)
	}
	l.logger.Log("StartRecordingInFile", "end")
	return
}
func (l *loggerMiddleware) StopRecoding(ctx context.Context, recoderIP, deviceName string) (err error) {
	l.logger.Log(
		"StopRecoding", "start",
		"recoderIP", recoderIP,
		"deviceName", deviceName,
	)
	if err = l.svc.StopRecoding(ctx, recoderIP, deviceName); err != nil {
		l.logger.Log(
			"StopRecoding", "err",
			"recoderIP", recoderIP,
			"deviceName", deviceName,
			"err", err,
		)
	}
	l.logger.Log("StopRecoding", "end")
	return
}

func (l *loggerMiddleware) StartRecordingOnPlayer(ctx context.Context, playerIP, recoderIP, deviceName string, channels, rate int) (err error) {
	l.logger.Log(
		"StopRecoding", "start",
		"playerIP", playerIP,
		"recoderIP", recoderIP,
		"deviceName", deviceName,
		"channels", channels,
		"rate", rate,
	)
	if err = l.svc.StartRecordingOnPlayer(ctx, playerIP, recoderIP, deviceName, channels, rate); err != nil {
		l.logger.Log(
			"StopRecoding", "start",
			"playerIP", playerIP,
			"recoderIP", recoderIP,
			"deviceName", deviceName,
			"channels", channels,
			"err", err,
		)
	}
	l.logger.Log("StopRecoding", "end")
	return
}

// NewLoggerMiddleware logger for server
func NewLoggerMiddleware(
	svc Server,
	logger log.Logger,
) Server {
	return &loggerMiddleware{
		svc:    svc,
		logger: logger,
	}
}
