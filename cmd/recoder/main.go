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
	"github.com/geoirb/sound-ethernet-streaming/pkg/converter"
	"github.com/geoirb/sound-ethernet-streaming/pkg/recorder"
	udp "github.com/geoirb/sound-ethernet-streaming/pkg/udp"
)

type configuration struct {
	Port        string `envconfig:"PORT" default:"8082"`
	UDPBuffSize int    `envconfig:"UDP_BUFF_SIZE" default:"64"`
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

	converter := converter.NewConverter()
	capture := capture.NewCapture(converter, cfg.UDPBuffSize)
	r5r := recorder.NewRecorder(
		udp,
		capture,
	)
	r5r = recorder.NewLoggerMiddleware(logger, r5r)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		level.Error(logger).Log("msg", "failed to turn up tcp connection", "err", err)
		os.Exit(1)
	}

	server := grpc.NewServer()
	recorder.RegisterRecorderServer(server, r5r)

	go server.Serve(lis)
	level.Error(logger).Log("msg", "recorder start", "port", cfg.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
}
