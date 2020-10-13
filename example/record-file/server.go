package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/geoirb/sound-server/pkg/recorder"
	"github.com/geoirb/sound-server/pkg/server"
	"github.com/geoirb/sound-server/pkg/udp"
	"github.com/geoirb/sound-server/pkg/wav"
)

type configuration struct {
	ServerIP string `envconfig:"SERVER_IP" default:"127.0.0.1"`

	PlayerPort   string `envconfig:"PLAYER_PORT" default:"8081"`
	RecorderPort string `envconfig:"RECODER_PORT" default:"8082"`

	UDPBuffSize int `envconfig:"UDP_BUF_SIZE" default:"1024"`

	AddrLayout   string `envconfig:"ADDRESS_LAYOUT" default:"%s:%s"`
	DeviceLayout string `envconfig:"DEVICE_LAYOUT" default:"%s:%s"`

	RecodeFile string `envconfig:"FILE" default:"/home/geo/go/src/github.com/geoirb/sound-server/audio/testRecode.wav"`
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
	recorder := recorder.NewClient(
		cfg.AddrLayout,
		cfg.RecorderPort,
	)
	udp := udp.NewUDP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		recorder,
		nil,
		udp,

		cfg.ServerIP,
		cfg.AddrLayout,
		cfg.DeviceLayout,
	)

	recorderIP := "127.0.0.1"
	receivePort := "8083"
	recorderDevice := "hw:0,0"
	file := "/home/geo/go/src/github.com/geoirb/sound-server/example/record-file/test.wav"
	svc = server.NewLoggerMiddleware(svc, logger)
	svc.StartFileRecoding(context.Background(), recorderIP, recorderDevice, 2, 44100, receivePort, file)
	level.Info(logger).Log("msg", "server start")

	time.Sleep(30 * time.Second)

	svc.StopFileRecoding(context.Background(), recorderIP, recorderDevice, receivePort)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)

}
