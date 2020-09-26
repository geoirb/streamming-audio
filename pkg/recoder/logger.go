package recoder

import (
	"context"

	"github.com/go-kit/kit/log"
)

// LoggerMiddleware ...
type LoggerMiddleware struct {
	logger  log.Logger
	recoder RecoderServer
}

// StartRecode log
func (l *LoggerMiddleware) StartRecode(ctx context.Context, in *StartRecodeRequest) (out *StartRecodeResponse, err error) {
	if out, err = l.recoder.StartRecode(ctx, in); err != nil {
		l.logger.Log(
			"in: %s\n err: %s",
			in.String(),
			err,
		)
	}
	return
}

// StopRecode ...
func (l *LoggerMiddleware) StopRecode(ctx context.Context, in *StopRecodeRequest) (out *StopRecodeResponse, err error) {
	if out, err = l.recoder.StopRecode(ctx, in); err != nil {
		l.logger.Log(
			"in: %s\n err: %s",
			in.String(),
			err,
		)
	}
	return
}

// NewLoggerMiddleware ...
func NewLoggerMiddleware(
	logger log.Logger,
	recoder RecoderServer,
) RecoderServer {
	return &LoggerMiddleware{
		logger:  logger,
		recoder: recoder,
	}
}
