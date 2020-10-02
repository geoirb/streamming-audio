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

	"github.com/geoirb/sound-ethernet-streaming/pkg/recorder"
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

	RecodeFile string `envconfig:"FILE" default:"/home/geo/go/src/github.com/geoirb/sound-ethernet-streaming/audio/testRecode.wav"`
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
	recorder := recorder.NewClient(
		cfg.HostLayout,
		cfg.RecorderPort,
	)
	udp := udp.NewUDP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		recorder,
		nil,
		udp,

		cfg.HostLayout,
		cfg.DeviceLayout,
	)

	recorderIP := "127.0.0.1"
	receivePort := "8083"
	recorderDevice := "hw:0,0"
	pwd, _ := os.Getwd()
	file := fmt.Sprintf("%s/%s", pwd, "example/play-file/server/test.wav")
	svc = server.NewLoggerMiddleware(svc, logger)
	svc.RecordingInFile(context.Background(), file, receivePort, recorderIP, recorderDevice, 2, 44100)
	level.Error(logger).Log("msg", "server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)

}
