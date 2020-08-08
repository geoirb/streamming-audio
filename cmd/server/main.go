package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/GeoIrb/sound-ethernet-streaming/pkg/server"
	udp "github.com/GeoIrb/sound-ethernet-streaming/pkg/udp/server"
	"github.com/GeoIrb/sound-ethernet-streaming/pkg/wav"
)

const (
	dstAddr = "127.0.0.1:8080"
	file    = "/home/geo/go/src/github.com/GeoIrb/sound-ethernet-streaming/audio/test.wav"
)

func main() {
	var (
		err  error
		data []byte
	)
	udpSrv := udp.NewServerUDP(dstAddr)
	if err = udpSrv.TurnOn(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer udpSrv.Shutdown()

	if data, err = ioutil.ReadFile(file); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	audio := wav.NewWAV()
	if err = audio.Parse(data); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s4v := server.NewServer()
	s4v.AddStreaming(udpSrv, audio)

	ctx, cancel := context.WithCancel(context.Background())
	s4v.Start(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	fmt.Printf("received signal, exiting signal %v\n", sig)
	cancel()
}
