package recorder

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	logger   log.Logger
	recorder RecorderServer
}

// Start log
func (l *loggerMiddleware) Start(ctx context.Context, in *StartSendRequest) (out *StartSendResponse, err error) {
	l.logger.Log("Start", "start", "in", in.String())
	if out, err = l.recorder.Start(ctx, in); err != nil {
		l.logger.Log("Start", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// Stop ...
func (l *loggerMiddleware) Stop(ctx context.Context, in *StopSendRequest) (out *StopSendResponse, err error) {
	l.logger.Log("Stop", "start", "in", in.String())
	if out, err = l.recorder.Stop(ctx, in); err != nil {
		l.logger.Log("Stop", "err", "in", in.String(), "err", err.Error())
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
