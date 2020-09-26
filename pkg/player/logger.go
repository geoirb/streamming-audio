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

// StartPlay log
func (l *LoggerMiddleware) StartPlay(ctx context.Context, in *StartPlayRequest) (out *StartPlayResponse, err error) {
	if out, err = l.player.StartPlay(ctx, in); err != nil {
		l.logger.Log(
			"in: %s\n err: %s",
			in.String(),
			err,
		)
	}
	return
}

// StopPlay log
func (l *LoggerMiddleware) StopPlay(ctx context.Context, in *StopPlayRequest) (out *StopPlayResponse, err error) {
	if out, err = l.player.StopPlay(ctx, in); err != nil {
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
	player PlayerServer,
) PlayerServer {
	return &LoggerMiddleware{
		logger: logger,
		player: player,
	}
}
