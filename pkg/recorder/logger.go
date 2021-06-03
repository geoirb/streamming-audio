package recorder

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	logger log.Logger
	server RecorderServer
}

// State logger
func (l *loggerMiddleware) State(ctx context.Context, in *StateRequest) (out *StateResponse, err error) {
	l.logger.Log("ReceiveStart", "start", "in", in.String())
	if out, err = l.server.State(ctx, in); err != nil {
		l.logger.Log("ReceiveStart", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// Start log
func (l *loggerMiddleware) Start(ctx context.Context, in *StartSendRequest) (out *StartSendResponse, err error) {
	l.logger.Log("Start", "start", "in", in.String())
	if out, err = l.server.Start(ctx, in); err != nil {
		l.logger.Log("Start", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// Stop log
func (l *loggerMiddleware) Stop(ctx context.Context, in *StopSendRequest) (out *StopSendResponse, err error) {
	l.logger.Log("Stop", "start", "in", in.String())
	if out, err = l.server.Stop(ctx, in); err != nil {
		l.logger.Log("Stop", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// NewLoggerMiddleware recoder
func NewLoggerMiddleware(
	logger log.Logger,
	recorder RecorderServer,
) RecorderServer {
	return &loggerMiddleware{
		logger: logger,
		server: recorder,
	}
}
