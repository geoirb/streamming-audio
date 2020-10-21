package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/geoirb/audio-service/pkg/player"
	"github.com/geoirb/audio-service/pkg/recorder"
	"github.com/geoirb/audio-service/pkg/server"
	"github.com/geoirb/audio-service/pkg/server/httpserver"
	"github.com/geoirb/audio-service/pkg/tcp"
	"github.com/geoirb/audio-service/pkg/wav"
)

type configuration struct {
	Port     string `envconfig:"PORT" default:"8000"`
	ServerIP string `envconfig:"SERVER_IP" default:"127.0.0.1"`

	PlayerPort   string `envconfig:"PLAYER_PORT" default:"8080"`
	RecorderPort string `envconfig:"RECODER_PORT" default:"8080"`

	UDPBuffSize int `envconfig:"UDP_BUF_SIZE" default:"1024"`

	AddrLayout   string `envconfig:"ADDRESS_LAYOUT" default:"%s:%s"`
	DeviceLayout string `envconfig:"DEVICE_LAYOUT" default:"%s:%s"`
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

	wav := wav.NewWAV()
	player := player.NewClient(
		cfg.AddrLayout,
		cfg.PlayerPort,
	)
	recorder := recorder.NewClient(
		cfg.AddrLayout,
		cfg.RecorderPort,
	)
	tcp := tcp.NewTCP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		recorder,
		player,
		tcp,

		cfg.ServerIP,
		cfg.AddrLayout,
		cfg.DeviceLayout,
	)
	svc = server.NewLoggerMiddleware(svc, logger)

	server := httpserver.NewServer(svc)

	go func() {
		level.Info(logger).Log("msg", "start server", "port", cfg.Port)
		if err := server.ListenAndServe(":" + cfg.Port); err != nil {
			level.Error(logger).Log("server run failure error %s", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)


	if err := server.Shutdown(); err != nil {
		level.Error(logger).Log("server shutdown failure %v", err)
	}
}
