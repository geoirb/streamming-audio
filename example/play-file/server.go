package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/geoirb/audio-service/pkg/player"
	"github.com/geoirb/audio-service/pkg/server"
	"github.com/geoirb/audio-service/pkg/tcp"
	"github.com/geoirb/audio-service/pkg/wav"
)

type configuration struct {
	ServerIP string `envconfig:"SERVER_IP" default:"127.0.0.1"`

	PlayerPort   string `envconfig:"PLAYER_PORT" default:"8080"`
	RecorderPort string `envconfig:"RECODER_PORT" default:"8080"`

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
	tcp := tcp.NewTCP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		nil,
		player,
		tcp,

		cfg.ServerIP,
		cfg.AddrLayout,
		cfg.DeviceLayout,
	)
	svc = server.NewLoggerMiddleware(svc, logger)
	level.Info(logger).Log("msg", "server start")

	pwd, _ := os.Getwd()
	file := pwd + "/audio/test.wav"
	playerIP := "127.0.0.1"
	playerPort := "8083"
	playerDeviceName := "hw:1,0"
	uuid, _, _, _ := svc.FilePlay(context.Background(), file, playerIP, playerPort, playerDeviceName)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)

	svc.FileStop(context.Background(), playerIP, playerPort, playerDeviceName, uuid)
}
