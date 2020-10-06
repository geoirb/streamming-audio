package player

import (
	"context"

	"github.com/go-kit/kit/log"
)

type loggerMiddleware struct {
	logger log.Logger
	player PlayerServer
}

// ReceiveStart log
func (l *loggerMiddleware) ReceiveStart(ctx context.Context, in *StartReceiveRequest) (out *StartReceiveResponse, err error) {
	l.logger.Log("ReceiveStart", "start", "in", in.String())
	if out, err = l.player.ReceiveStart(ctx, in); err != nil {
		l.logger.Log("ReceiveStart", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// ReceiveStop log
func (l *loggerMiddleware) ReceiveStop(ctx context.Context, in *StopReceiveRequest) (out *StopReceiveResponse, err error) {
	l.logger.Log("ReceiveStop", "start", "in", in.String())
	if out, err = l.player.ReceiveStop(ctx, in); err != nil {
		l.logger.Log("ReceiveStop", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// Play log
func (l *loggerMiddleware) Play(ctx context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	l.logger.Log("Play", "start", "in", in.String())
	if out, err = l.player.Play(ctx, in); err != nil {
		l.logger.Log("Play", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// Stop log
func (l *loggerMiddleware) Stop(ctx context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
	l.logger.Log("Stop", "start", "in", in.String())
	if out, err = l.player.Stop(ctx, in); err != nil {
		l.logger.Log("Stop", "err", "in", in.String(), "err", err.Error())
	}
	return
}

// ClearStorage log
func (l *loggerMiddleware) ClearStorage(ctx context.Context, in *ClearStorageRequest) (out *ClearStorageResponse, err error) {
	l.logger.Log("ClearStorage", "start", "in", in.String())
	if out, err = l.player.ClearStorage(ctx, in); err != nil {
		l.logger.Log("ClearStorage", "err", "in", in.String(), "err", err.Error())
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
