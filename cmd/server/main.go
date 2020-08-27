package main

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/geoirb/sound-ethernet-streaming/pkg/server"
	udp "github.com/geoirb/sound-ethernet-streaming/pkg/udp/server"
	"github.com/geoirb/sound-ethernet-streaming/pkg/wav"
)

type configuration struct {
	DstAddress string `envconfig:"DST_ADDRESS" default:"255.255.255.255:8080"`
	File       string `envconfig:"FILE" default:"/home/geo/go/src/github.com/geoirb/sound-ethernet-streaming/audio/test.wav"`
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	_ = level.Info(logger).Log("msg", "initializing")

	var (
		err  error
		data []byte
		cfg  configuration
	)
	if err = envconfig.Process("", &cfg); err != nil {
		_ = level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(1)
	}

	if data, err = ioutil.ReadFile(cfg.File); err != nil {
		_ = level.Error(logger).Log("msg", "failed to read file", "file", cfg.File, "err", err)
		os.Exit(1)
	}
	source := wav.NewWAV()
	if err = source.Parse(data); err != nil {
		_ = level.Error(logger).Log("msg", "failed to parse wav", "err", err)
		os.Exit(1)
	}

	udpSrv := udp.NewServerUDP(cfg.DstAddress)
	if err = udpSrv.TurnOn(); err != nil {
		_ = level.Error(logger).Log("msg", "failed to turn on udp server", "err", err)
		os.Exit(1)
	}
	defer udpSrv.Shutdown()

	s4v := server.NewServer()
	if err = s4v.AddStreaming(udpSrv, source); err != nil {
		_ = level.Error(logger).Log("msg", "add streaming", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s4v.Start(ctx)
	_ = level.Error(logger).Log("msg", "server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	_ = level.Error(logger).Log("msg", "received signal, exiting signal", "signal", sig)
	cancel()
}
