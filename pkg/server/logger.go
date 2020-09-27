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
	l.logger.Log("StartSendingFile", "start", "destIP", destIP, "destIP", destPort, "fileName", fileName)
	if storageUUID, channels, rate, err = l.svc.StartSendingFile(ctx, destIP, destPort, fileName); err != nil {
		l.logger.Log("StartSendingFile", "err", "destIP", destIP, "destPort", destPort, "fileName", fileName, "err", err)
	}
	l.logger.Log("StartSendingFile", "end")
	return
}

func (l *loggerMiddleware) StopSending(ctx context.Context, destIP, destPort string) (err error) {
	l.logger.Log("StartSendingFile", "start", "destIP", destIP, "destIP", destPort)
	if err = l.svc.StopSending(ctx, destIP, destPort); err != nil {
		l.logger.Log("StartSendingFile", "err", "destIP", destIP, "destPort", destPort, "err", err)
	}
	l.logger.Log("StopSending", "end")
	return
}

func (l *loggerMiddleware) StartPlaying(ctx context.Context, playerIP, deviceName, storageUUID string, channels uint16, rate uint32) (err error) {
	l.logger.Log("StartSendingFile", "start", "destIP", playerIP, "deviceName", "storageUUID", storageUUID, deviceName, "channels", channels, "rate", rate)
	if err = l.svc.StartPlaying(ctx, playerIP, deviceName, storageUUID, channels, rate); err != nil {
		l.logger.Log("StartSendingFile", "err", "destIP", playerIP, "deviceName", "storageUUID", storageUUID, deviceName, "channels", channels, "rate", rate, "err", err)
	}
	l.logger.Log("StartPlaying", "end")
	return
}

func (l *loggerMiddleware) StopPlaying(ctx context.Context, playerIP, deviceName string) (err error) {
	l.logger.Log("StartSendingFile", "start", "destIP", playerIP, "deviceName", deviceName)
	if err = l.svc.StopPlaying(ctx, playerIP, deviceName); err != nil {
		l.logger.Log("StartSendingFile", "err", "destIP", playerIP, "deviceName", deviceName, "err", err)
	}
	l.logger.Log("StopPlaying", "end")
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
