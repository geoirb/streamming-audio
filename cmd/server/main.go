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

	"github.com/geoirb/sound-ethernet-streaming/pkg/controller/media"
	"github.com/geoirb/sound-ethernet-streaming/pkg/server"
	udp "github.com/geoirb/sound-ethernet-streaming/pkg/udp/server"
	"github.com/geoirb/sound-ethernet-streaming/pkg/wav"
)

type configuration struct {
	Port        string `envconfig:"PORT" default:"8081"`
	UDPBuffSize int    `envconfig:"UDP_BUF_SIZE" default:"1024"`
	HostLayout  string `envconfig:"HOST_LAYOUT" default:"%s:%s"`
	File        string `envconfig:"FILE" default:"/home/geo/go/src/github.com/geoirb/sound-ethernet-streaming/audio/test.wav"`
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

	w1v := wav.NewWAV()
	m3a := media.NewMediaController(cfg.HostLayout, cfg.Port)
	u1p := udp.NewServer(cfg.UDPBuffSize)
	s4r := server.NewServer(
		cfg.HostLayout,
		w1v,
		m3a,
		u1p,
	)
	fmt.Println(s4r.AddFileMedia(context.Background(), "127.0.0.1", "8082", "hw:0,0", cfg.File))

	level.Error(logger).Log("msg", "server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
}
