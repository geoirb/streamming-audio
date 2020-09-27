package player

import (
	"context"

	"github.com/go-kit/kit/log"
)

// LoggerMiddleware ...
type LoggerMiddleware struct {
	logger log.Logger
	player PlayerServer
}

// StartReceive log
func (l *LoggerMiddleware) StartReceive(ctx context.Context, in *StartReceiveRequest) (out *StartReceiveRequest, err error) {
	if out, err = l.player.StartReceive(ctx, in); err != nil {
		l.logger.Log("function", "StartReceive", "in", in.String(), "err", err.Error())
	}
	return
}

// StopReceive log
func (l *LoggerMiddleware) StopReceive(ctx context.Context, in *StopReceiveRequest) (out *StopReceiveRequest, err error) {
	if out, err = l.player.StopReceive(ctx, in); err != nil {
		l.logger.Log("function", "StopReceive", "in", in.String(), "err", err.Error())
	}
	return
}

// StartPlay log
func (l *LoggerMiddleware) StartPlay(ctx context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	if out, err = l.player.StartPlay(ctx, in); err != nil {
		l.logger.Log("function", "StartPlay", "in", in.String(), "err", err.Error())
	}
	return
}

// StopPlay log
func (l *LoggerMiddleware) StopPlay(ctx context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
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
	return &LoggerMiddleware{
		logger: logger,
		player: player,
	}
}
