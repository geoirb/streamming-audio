package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/geoirb/sound-ethernet-streaming/pkg/server"
	"github.com/geoirb/sound-ethernet-streaming/pkg/wav"
)

type configuration struct {
	HostLayout string `envconfig:"HOST_LAYOUT" default:"%s:%s"`
	DstAddress string `envconfig:"DST_ADDRESS" default:"255.255.255.255:8080"`
	File       string `envconfig:"FILE" default:"/home/geo/go/src/github.com/geoirb/sound-ethernet-streaming/audio/test.wav"`
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	_ = level.Info(logger).Log("msg", "initializing")

	var (
		err error
		cfg configuration
	)
	if err = envconfig.Process("", &cfg); err != nil {
		_ = level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(1)
	}

	w1v := wav.NewWAV()
	s4v := server.NewServer(
		cfg.HostLayout,
		w1v,
	)

	_ = level.Error(logger).Log("msg", "server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	_ = level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
}
