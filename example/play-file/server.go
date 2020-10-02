package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/geoirb/sound-ethernet-streaming/pkg/player"
	"github.com/geoirb/sound-ethernet-streaming/pkg/server"
	"github.com/geoirb/sound-ethernet-streaming/pkg/udp"
	"github.com/geoirb/sound-ethernet-streaming/pkg/wav"
)

type configuration struct {
	PlayerPort   string `envconfig:"PLAYER_PORT" default:"8081"`
	RecorderPort string `envconfig:"RECODER_PORT" default:"8082"`

	UDPBuffSize int `envconfig:"UDP_BUF_SIZE" default:"1024"`

	HostLayout   string `envconfig:"HOST_LAYOUT" default:"%s:%s"`
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
		cfg.HostLayout,
		cfg.PlayerPort,
	)
	udp := udp.NewUDP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		nil,
		player,
		udp,

		cfg.HostLayout,
		cfg.DeviceLayout,
	)
	svc = server.NewLoggerMiddleware(svc, logger)
	level.Error(logger).Log("msg", "server start")

	pwd, _ := os.Getwd()
	file := fmt.Sprintf("%s/%s", pwd, "example/play-file/server/test.wav")
	playerIP := "127.0.0.1"
	playerPort := "8083"
	playerDeviceName := "hw:0,0"
	uuid, _, _, _ := svc.PlayAudioFile(context.Background(), playerIP, playerPort, file, playerDeviceName)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)

	svc.Stop(context.Background(), playerIP, playerPort, playerDeviceName, uuid)
}
