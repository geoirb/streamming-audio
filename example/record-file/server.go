package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"github.com/geoirb/audio-service/pkg/recorder"
	"github.com/geoirb/audio-service/pkg/server"
	"github.com/geoirb/audio-service/pkg/tcp"
	"github.com/geoirb/audio-service/pkg/wav"
)

type configuration struct {
	ServerIP string `envconfig:"SERVER_IP" default:"127.0.0.1"`

	PlayerPort   string `envconfig:"PLAYER_PORT" default:"8081"`
	RecorderPort string `envconfig:"RECODER_PORT" default:"8082"`

	UDPBuffSize int `envconfig:"UDP_BUF_SIZE" default:"1024"`

	AddrLayout   string `envconfig:"ADDRESS_LAYOUT" default:"%s:%s"`
	DeviceLayout string `envconfig:"DEVICE_LAYOUT" default:"%s:%s"`

	RecodeFile string `envconfig:"FILE" default:"/home/geo/go/src/github.com/geoirb/audio-service/audio/testRecode.wav"`
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
	tcp := tcp.NewTCP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		recorder,
		nil,
		tcp,

		cfg.ServerIP,
		cfg.AddrLayout,
		cfg.DeviceLayout,
	)

	recorderIP := "127.0.0.1"
	receivePort := "8083"
	recorderDevice := "hw:0,0"
	pwd, _ := os.Getwd()
	file := pwd + "/example/record-file/test.wav"
	svc = server.NewLoggerMiddleware(svc, logger)
	svc.StartFileRecording(context.Background(), recorderIP, recorderDevice, 2, 44100, receivePort, file)
	level.Info(logger).Log("msg", "server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)

	svc.StopFileRecording(context.Background(), recorderIP, recorderDevice, receivePort)

}
