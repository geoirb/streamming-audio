package player

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	logger log.Logger
	player PlayerServer
}

// StartReceive log
func (l *loggerMiddleware) StartReceive(ctx context.Context, in *StartReceiveRequest) (out *StartReceiveResponse, err error) {
	l.logger.Log("function", "StartReceive", "in", in.String())
	if out, err = l.player.StartReceive(ctx, in); err != nil {
		l.logger.Log("function", "StartReceive", "in", in.String(), "err", err.Error())
	}
	return
}

// StopReceive log
func (l *loggerMiddleware) StopReceive(ctx context.Context, in *StopReceiveRequest) (out *StopReceiveResponse, err error) {
	l.logger.Log("function", "StopReceive", "in", in.String())
	if out, err = l.player.StopReceive(ctx, in); err != nil {
		l.logger.Log("function", "StopReceive", "in", in.String(), "err", err.Error())
	}
	return
}

// StartPlay log
func (l *loggerMiddleware) StartPlay(ctx context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	l.logger.Log("function", "StartPlay", "in", in.String())
	if out, err = l.player.StartPlay(ctx, in); err != nil {
		l.logger.Log("function", "StartPlay", "in", in.String(), "err", err.Error())
	}
	return
}

// StopPlay log
func (l *loggerMiddleware) StopPlay(ctx context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
	l.logger.Log("function", "StopPlay", "in", in.String())
	if out, err = l.player.StopPlay(ctx, in); err != nil {
		l.logger.Log("function", "StopPlay", "in", in.String(), "err", err.Error())
	}
	return
}

// NewLoggerMiddleware ...
func NewLoggerMiddleware(
	logger log.Logger,
	player PlayerServer,
) PlayerServer {
	return &loggerMiddleware{
		logger: logger,
		player: player,
	}
}
