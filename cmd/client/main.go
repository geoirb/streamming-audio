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

	PlaybackDeviceName string `envconfig:"PLAYBACK_DEVICE_NAME" default:"hw:0,0"`
	PlaybackChannels   int    `envconfig:"PLAYBACK_DEVICE_NAME" default:"2"`
	PlaybackRate       int    `envconfig:"PLAYBACK_DEVICE_NAME" default:"44100"`
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

	udpClt := udp.NewClientUDP(cfg.Port, cfg.UDPBufSize)
	if err = udpClt.Connect(); err != nil {
		_ = level.Error(logger).Log("msg", "failed to connect udp server", "err", err)
		os.Exit(1)
	}
	defer udpClt.Disconnect()

	p6k := playback.NewPlayback(
		cfg.PlaybackDeviceName,
		cfg.PlaybackChannels,
		cfg.PlaybackRate,
	)
	if err = p6k.Device(); err != nil {
		_ = level.Error(logger).Log("msg", "failed to connect to playback device", "err", err)
		os.Exit(1)
	}

	c7r := converter.NewConverter()
	c2h := cash.NewCash()

	m := client.NewClient(c7r)

	if err = m.Add(p6k, udpClt, c2h); err != nil {
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
	p6k.Disconnect()
	cancel()
}
