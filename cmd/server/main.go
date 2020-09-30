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

	"github.com/geoirb/sound-ethernet-streaming/pkg/player"
	"github.com/geoirb/sound-ethernet-streaming/pkg/recorder"
	"github.com/geoirb/sound-ethernet-streaming/pkg/server"
	"github.com/geoirb/sound-ethernet-streaming/pkg/udp"
	"github.com/geoirb/sound-ethernet-streaming/pkg/wav"
)

type configuration struct {
	PlayerPort   string `envconfig:"PLAYER_PORT" default:"8081"`
	RecorderPort string `envconfig:"RECODER_PORT" default:"8082"`

	UDPBuffSize int `envconfig:"UDP_BUF_SIZE" default:"1024"`

	HostLayout   string `envconfig:"HOST_LAYOUT" default:"%s:%s"`
	DeviceLayout string `envconfig:"DEVICE_LAYOUT" default:"%s:%s"`

	PlayFile   string `envconfig:"FILE" default:"/home/geo/go/src/github.com/geoirb/sound-ethernet-streaming/audio/test.wav"`
	RecodeFile string `envconfig:"FILE" default:"/home/geo/go/src/github.com/geoirb/sound-ethernet-streaming/audio/testRecode.wav"`
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
		cfg.HostLayout,
		cfg.PlayerPort,
	)
	recorder := recorder.NewClient(
		cfg.HostLayout,
		cfg.RecorderPort,
	)
	udp := udp.NewUDP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		recorder,
		player,
		udp,

		cfg.HostLayout,
		cfg.DeviceLayout,
	)
	svc = server.NewLoggerMiddleware(svc, logger)
	svc.RecordingInFile(context.Background(), cfg.RecodeFile, "8083", "127.0.0.1", "hw:0,0", 2, 44100)
	time.Sleep(time.Second * 5)
	svc.StopRecoding(context.Background(), "127.0.0.1", "hw:0,0")

	level.Error(logger).Log("msg", "server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Error(logger).Log("msg", "received signal, exiting signal", "signal", <-c)

}

// uuid, channels, rate, _ := svc.PlayAudioFile(context.Background(), "127.0.0.1", "8083", cfg.PlayFile, "hw:1,0")
// time.Sleep(5 * time.Second)
// svc.Pause(context.Background(), "127.0.0.1", "hw:1,0")
// time.Sleep(10 * time.Second)
// svc.Play(context.Background(), "127.0.0.1", uuid, "hw:1,0", channels, rate)
// time.Sleep(5 * time.Second)
// svc.Stop(context.Background(), "127.0.0.1", "8083", "hw:1,0", uuid)
// svc.Play(context.Background(), "127.0.0.1", uuid, "hw:1,0", channels, rate)
