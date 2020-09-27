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

	"github.com/geoirb/sound-ethernet-streaming/pkg/converter"
	"github.com/geoirb/sound-ethernet-streaming/pkg/playback"
	"github.com/geoirb/sound-ethernet-streaming/pkg/player"
	"github.com/geoirb/sound-ethernet-streaming/pkg/storage"
	udp "github.com/geoirb/sound-ethernet-streaming/pkg/udp"
)

type configuration struct {
	Port        string `envconfig:"PORT" default:"8081"`
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

	converter := converter.NewConverter()
	playback := playback.NewPlayback(
		converter,
		cfg.UDPBuffSize,
	)

	storage := storage.NewStorage()

	p4r := player.NewPlayer(
		udp,
		playback,
		storage,
	)
	p4r = player.NewLoggerMiddleware(logger, p4r)

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		level.Error(logger).Log("msg", "failed to turn up tcp connection", "err", err)
		os.Exit(1)
	}

	server := grpc.NewServer()
	player.RegisterPlayerServer(server, p4r)

	go server.Serve(lis)
	level.Error(logger).Log("msg", "player start", "port", cfg.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)
}
