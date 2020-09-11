package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	controller "github.com/geoirb/sound-ethernet-streaming/pkg/controller/media"
	"github.com/geoirb/sound-ethernet-streaming/pkg/converter"
	"github.com/geoirb/sound-ethernet-streaming/pkg/media"
	"github.com/geoirb/sound-ethernet-streaming/pkg/playback"
	udp "github.com/geoirb/sound-ethernet-streaming/pkg/udp/client"
)

type configuration struct {
	Port       string `envconfig:"PORT" default:"8081"`
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

	udpClt := udp.NewClient(cfg.UDPBufSize)

	c7r := converter.NewConverter()
	p6k := playback.NewPlayback(c7r)

	m3a := media.NewMedia(
		udpClt,
		p6k,
	)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed to turn up tcp connection", "err", err)
		os.Exit(1)
	}

	server := grpc.NewServer()
	controller.RegisterMediaServer(server, m3a)

	_ = level.Error(logger).Log("msg", "server start", "port", cfg.Port)
	server.Serve(lis)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	_ = level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
}
