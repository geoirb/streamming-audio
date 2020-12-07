package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"audio-service/pkg/player"
	"audio-service/pkg/recorder"
	"audio-service/pkg/server"
	"audio-service/pkg/tcp"
	"audio-service/pkg/wav"
)

type configuration struct {
	ServerIP string `envconfig:"SERVER_IP" default:"127.0.0.1"`

	PlayerPort   string `envconfig:"PLAYER_PORT" default:"8081"`
	RecorderPort string `envconfig:"RECODER_PORT" default:"8082"`

	UDPBuffSize int `envconfig:"UDP_BUF_SIZE" default:"1024"`

	AddrLayout   string `envconfig:"ADDRESS_LAYOUT" default:"%s:%s"`
	DeviceLayout string `envconfig:"DEVICE_LAYOUT" default:"%s:%s"`
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	level.Info(logger).Log("msg", "initializing")

	var (
		err error
		cfg configuration
	)
	if err = envconfig.Process("", &cfg); err != nil {
		level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(1)
	}

	wav := wav.NewWAV()
	player := player.NewClient(
		cfg.AddrLayout,
		cfg.PlayerPort,
	)
	recorder := recorder.NewClient(
		cfg.AddrLayout,
		cfg.RecorderPort,
	)
	tcp := tcp.NewTCP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		recorder,
		player,
		tcp,

		cfg.ServerIP,
		cfg.AddrLayout,
		cfg.DeviceLayout,
	)
	svc = server.NewLoggerMiddleware(svc, logger)
	uuid, _ := svc.PlayFromRecorder(context.Background(), "127.0.0.1", "8083", "hw:1,0", 2, 44100, "127.0.0.1", "hw:0,0")
	level.Info(logger).Log("msg", "server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
	svc.StopFromRecorder(context.Background(), "127.0.0.1", "8083", "hw:0,0", uuid, "127.0.0.1", "hw:0,0")
}
