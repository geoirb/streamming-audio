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

	"github.com/geoirb/sound-ethernet-streaming/pkg/capture"
	"github.com/geoirb/sound-ethernet-streaming/pkg/recoder"
	udp "github.com/geoirb/sound-ethernet-streaming/pkg/udp"
)

type configuration struct {
	Port        string `envconfig:"PORT" default:"8082"`
	UDPBuffSize int    `envconfig:"UDP_BUFF_SIZE" default:"1024"`
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

	udp := udp.NewUDP(cfg.UDPBuffSize)

	capture := capture.NewCapture()
	r5r := recoder.NewRecoder(
		udp,
		capture,
	)
	r5r = recoder.NewLoggerMiddleware(logger, r5r)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		level.Error(logger).Log("msg", "failed to turn up tcp connection", "err", err)
		os.Exit(1)
	}

	server := grpc.NewServer()
	recoder.RegisterRecoderServer(server, r5r)

	level.Error(logger).Log("msg", "server start", "port", cfg.Port)
	server.Serve(lis)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
}
