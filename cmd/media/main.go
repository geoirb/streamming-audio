package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/geoirb/sound-ethernet-streaming/pkg/cash"
	"github.com/geoirb/sound-ethernet-streaming/pkg/client"
	"github.com/geoirb/sound-ethernet-streaming/pkg/converter"
	"github.com/geoirb/sound-ethernet-streaming/pkg/playback"
	udp "github.com/geoirb/sound-ethernet-streaming/pkg/udp/client"
)

type configuration struct {
	Port       string `envconfig:"PORT" default:":8080"`
	UDPBufSize int    `envconfig:"UDP_BUF_SIZE" default:"1024"`

	PlaybackDeviceName string `envconfig:"PLAYBACK_DEVICE_NAME" default:"hw:1,0"`
	PlaybackChannels   int    `envconfig:"PLAYBACK_CHANELS" default:"2"`
	PlaybackRate       int    `envconfig:"PLAYBACK_RATE" default:"44100"`
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

	udpClt, err := udp.NewClientUDP(cfg.Port, cfg.UDPBufSize)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to connect udp server", "err", err)
		os.Exit(1)
	}
	defer udpClt.Disconnect()

	c7r := converter.NewConverter()
	p6k, err := playback.NewPlayback(
		cfg.PlaybackDeviceName,
		cfg.PlaybackChannels,
		cfg.PlaybackRate,
		c7r,
	)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to connect to playback device", "err", err)
		os.Exit(1)
	}
	defer p6k.Disconnect()

	c2h := cash.NewCash()
	c := client.NewClient()

	if err = c.Add(p6k, udpClt, c2h); err != nil {
		_ = level.Error(logger).Log("msg", "failed to add in client", "err", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.Start(ctx)
	_ = level.Error(logger).Log("msg", "server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	_ = level.Error(logger).Log("msg", "received signal, exiting signal", "signal", sig)
	cancel()
}
