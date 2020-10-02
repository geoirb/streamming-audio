package recorder

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	logger   log.Logger
	recorder RecorderServer
}

// StartRecord log
func (l *loggerMiddleware) StartRecord(ctx context.Context, in *StartRecordRequest) (out *StartRecordResponse, err error) {
	l.logger.Log("StartRecord", "start", "in", in.String())
	if out, err = l.recorder.StartRecord(ctx, in); err != nil {
		l.logger.Log("StartRecord", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// StopRecord ...
func (l *loggerMiddleware) StopRecord(ctx context.Context, in *StopRecordRequest) (out *StopRecordResponse, err error) {
	l.logger.Log("StopRecord", "start", "in", in.String())
	if out, err = l.recorder.StopRecord(ctx, in); err != nil {
		l.logger.Log("StopRecord", "err", "in", in.String(), "err", err.Error())
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
