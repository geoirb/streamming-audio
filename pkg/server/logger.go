package server

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	svc    Server
	logger log.Logger
}

func (l *loggerMiddleware) PlayAudioFile(ctx context.Context, playerIP, playerPort, fileName, deviceName string) (storageUUID string, channels uint16, rate uint32, err error) {
	l.logger.Log(
		"PlayAudioFile", "start",
		"playerIP", playerIP,
		"playerPort", playerPort,
		"fileName", fileName,
		"deviceName", deviceName,
	)
	if storageUUID, channels, rate, err = l.svc.PlayAudioFile(ctx, playerIP, playerPort, fileName, deviceName); err != nil {
		l.logger.Log(
			"PlayAudioFile", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"fileName", fileName,
			"deviceName", deviceName,
			"err", err,
		)
	}
	l.logger.Log("PlayAudioFile", "end")
	return
}

func (l *loggerMiddleware) Play(ctx context.Context, playerIP, storageUUID, deviceName string, channels uint16, rate uint32) (err error) {
	l.logger.Log(
		"Play", "start",
		"playerIP", playerIP,
		"storageUUID", storageUUID,
		"deviceName", deviceName,
		"channels", channels,
		"rate", rate,
	)
	if err = l.svc.Play(ctx, playerIP, storageUUID, deviceName, channels, rate); err != nil {
		l.logger.Log(
			"Play", "err",
			"playerIP", playerIP,
			"storageUUID", storageUUID,
			"deviceName", deviceName,
			"channels", channels,
			"rate", rate,
			"err", err,
		)
	}
	l.logger.Log("Play", "end")
	return
}

func (l *loggerMiddleware) Pause(ctx context.Context, playerIP, deviceName string) (err error) {
	l.logger.Log(
		"Pause", "start",
		"playerIP", playerIP,
		"deviceName", deviceName,
	)
	if err = l.svc.Pause(ctx, playerIP, deviceName); err != nil {
		l.logger.Log(
			"Pause", "err",
			"playerIP", playerIP,
			"deviceName", deviceName,
			"err", err,
		)
	}
	l.logger.Log("Pause", "end")
	return
}

func (l *loggerMiddleware) Stop(ctx context.Context, playerIP, playerPort, deviceName, storageUUID string) (err error) {
	l.logger.Log(
		"Stop", "start",
		"playerIP", playerIP,
		"playerPort", playerPort,
		"deviceName", deviceName,
		"storageUUID", storageUUID,
	)
	if err = l.svc.Stop(ctx, playerIP, playerPort, deviceName, storageUUID); err != nil {
		l.logger.Log(
			"Stop", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"deviceName", deviceName,
			"storageUUID", storageUUID,
			"err", err,
		)
	}
	l.logger.Log("Stop", "end")
	return
}

func (l *loggerMiddleware) RecordingOnPlayer(ctx context.Context, playerIP, playerPort, playerDeviceNameName, recoderIP, recorderDeviceName string, channels, rate int) (storageUUID string, err error) {
	l.logger.Log(
		"StartSendingOnPlayer", "start",
		"playerIP", playerIP,
		"playerPort", playerPort,
		"playerDeviceNameName", playerDeviceNameName,
		"recoderIP", recoderIP,
		"recorderDeviceName", recorderDeviceName,
		"channels", channels,
		"rate", rate,
	)
	if storageUUID, err = l.svc.RecordingOnPlayer(ctx, playerIP, playerPort, playerDeviceNameName, recoderIP, recorderDeviceName, channels, rate); err != nil {
		l.logger.Log(
			"StartSendingInFile", "err",
			"playerIP", playerIP,
			"playerPort", playerPort,
			"playerDeviceNameName", playerDeviceNameName,
			"recoderIP", recoderIP,
			"recorderDeviceName", recorderDeviceName,
			"channels", channels,
			"rate", rate,
			"err", err,
		)
	}
	l.logger.Log("StartSendingOnPlayer", "end")
	return
}

func (l *loggerMiddleware) RecordingInFile(c context.Context, fileName, receivePort, recoderIP, recoderDeviceName string, channels, rate int) (err error) {
	l.logger.Log(
		"RecordingInFile", "start",
		"fileName", fileName,
		"receivePort", receivePort,
		"recoderIP", recoderIP,
		"recoderIP", recoderIP,
		"recoderDeviceName", recoderDeviceName,
		"channels", channels,
		"rate", rate,
	)
	if err = l.svc.RecordingInFile(c, fileName, receivePort, recoderIP, recoderDeviceName, channels, rate); err != nil {
		l.logger.Log(
			"RecordingInFile", "err",
			"fileName", fileName,
			"receivePort", receivePort,
			"recoderIP", recoderIP,
			"recoderIP", recoderIP,
			"recoderDeviceName", recoderDeviceName,
			"channels", channels,
			"rate", rate,
			"err", err,
		)
	}
	l.logger.Log("RecordingInFile", "end")
	return
}

func (l *loggerMiddleware) StopRecoding(c context.Context, recoderIP, deviceName string) (err error) {
	l.logger.Log(
		"StopRecoding", "start",
		"recoderIP", recoderIP,
		"deviceName", deviceName,
	)
	if err = l.svc.StopRecoding(c, recoderIP, deviceName); err != nil {
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
