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

	"github.com/geoirb/ausio-service/pkg/capture"
	"github.com/geoirb/ausio-service/pkg/converter"
	"github.com/geoirb/ausio-service/pkg/recorder"
	tcp "github.com/geoirb/ausio-service/pkg/tcp"
)

type configuration struct {
	Port        string `envconfig:"PORT" default:"8080"`
	UDPBuffSize int    `envconfig:"UDP_BUFF_SIZE" default:"32"`
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

	tcp := tcp.NewTCP(cfg.UDPBuffSize)

	converter := converter.NewConverter()
	capture := capture.NewCapture(converter, cfg.UDPBuffSize)
	r5r := recorder.NewRecorder(
		tcp,
		capture,
	)
	r5r = recorder.NewLoggerMiddleware(logger, r5r)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		level.Error(logger).Log("msg", "failed to turn up tcp connection", "err", err)
		os.Exit(1)
	}
	defer lis.Close()

	server := grpc.NewServer()
	recorder.RegisterRecorderServer(server, r5r)

	go server.Serve(lis)
	level.Info(logger).Log("msg", "recorder start", "port", cfg.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
}
