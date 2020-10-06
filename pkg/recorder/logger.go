package recorder

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	logger   log.Logger
	recorder RecorderServer
}

// StartSend log
func (l *loggerMiddleware) StartSend(ctx context.Context, in *StartSendRequest) (out *StartSendResponse, err error) {
	l.logger.Log("StartSend", "start", "in", in.String())
	if out, err = l.recorder.StartSend(ctx, in); err != nil {
		l.logger.Log("StartSend", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// StopSend ...
func (l *loggerMiddleware) StopSend(ctx context.Context, in *StopSendRequest) (out *StopSendResponse, err error) {
	l.logger.Log("StopSend", "start", "in", in.String())
	if out, err = l.recorder.StopSend(ctx, in); err != nil {
		l.logger.Log("StopSend", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// NewLoggerMiddleware ...
func NewLoggerMiddleware(
	logger log.Logger,
	recorder RecorderServer,
) RecorderServer {
	return &loggerMiddleware{
		logger:   logger,
		recorder: recorder,
	}
}
