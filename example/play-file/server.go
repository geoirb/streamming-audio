package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"

	"audio-service/pkg/player"
	"audio-service/pkg/server"
	"audio-service/pkg/tcp"
	"audio-service/pkg/wav"
)

type configuration struct {
	ServerIP string `envconfig:"SERVER_IP" default:"127.0.0.1"`

	PlayerPort   string `envconfig:"PLAYER_PORT" default:"8080"`
	RecorderPort string `envconfig:"RECODER_PORT" default:"8080"`

	UDPBuffSize int `envconfig:"UDP_BUF_SIZE" default:"1024"`

	AddrLayout   string `envconfig:"ADDRESS_LAYOUT" default:"%s:%s"`
	DeviceLayout string `envconfig:"DEVICE_LAYOUT" default:"%s:%s"`
}

type playerInfo struct {
	Start  bool
	UUID   string
	IP     string
	Port   string
	Device string
	File   string
}

var playerConf map[string]playerInfo = map[string]playerInfo{
	"1": {
		IP:     "127.0.0.1",
		Port:   "8081",
		Device: "hw:0,0",
		File:   "/home/geo/go/src/audio-service",
	},
	"2": {
		IP:     "192.168.0.106",
		Port:   "8081",
		Device: "hw:0,0",
		File:   "/audio/NAME_TEST_FILE_ON_AUDIO_DIR.wav",
	},
	"3": {
		IP:     "192.168.0.106",
		Port:   "8081",
		Device: "hw:0,0",
		File:   "/audio/NAME_TEST_FILE_ON_AUDIO_DIR.wav",
	},
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
	tcp := tcp.NewTCP(cfg.UDPBuffSize)
	svc := server.NewServer(
		wav,
		nil,
		player,
		tcp,

		cfg.ServerIP,
		cfg.AddrLayout,
		cfg.DeviceLayout,
	)
	svc = server.NewLoggerMiddleware(svc, logger)
	level.Info(logger).Log("msg", "server start")

	// example
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Input number of player")
	for scanner.Scan() {
		num := scanner.Text()
		p, isExist := playerConf[num]
		if !isExist {
			fmt.Printf("player num %v not exist\n", num)
		}
		if !p.Start {
			if uuid, _, _, _, err := svc.FilePlay(context.Background(), p.File, p.IP, p.Port, p.Device); err == nil {
				p.UUID = uuid
				p.Start = true
				playerConf[num] = p
				continue
			} else {
				fmt.Println(err)
			}
		} else {
			if err := svc.FileStop(context.Background(), p.IP, p.Port, p.Device, p.UUID); err == nil {
				p.Start = false
				playerConf[num] = p
				continue
			} else {
				fmt.Println(err)
			}
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal, exiting signal", "signal", <-c)

}
