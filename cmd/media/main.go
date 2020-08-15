package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/GeoIrb/sound-ethernet-streaming/pkg/cash"
	"github.com/GeoIrb/sound-ethernet-streaming/pkg/converter"
	"github.com/GeoIrb/sound-ethernet-streaming/pkg/media"
	"github.com/GeoIrb/sound-ethernet-streaming/pkg/playback"
	udp "github.com/GeoIrb/sound-ethernet-streaming/pkg/udp/client"
)

const (
	sizeData  = 100
	localAddr = ":8080"

	deviceName = "hw:1,0"
	channels   = 2
	rate       = 44100
	buffSize   = 1024
)

func main() {
	udpClt := udp.NewClientUDP(localAddr, buffSize)
	udpClt.Connect()
	defer udpClt.Disconnect()

	p6k := playback.NewPlayback(
		deviceName,
		channels,
		rate,
	)
	p6k.Device()

	c7r := converter.NewConverter()
	c2h := cash.NewCash()

	m := media.NewMedia(c7r)

	m.Add(p6k, udpClt, c2h)
	ctx, cancel := context.WithCancel(context.Background())
	m.Start(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	fmt.Printf("received signal, exiting signal %v\n", sig)
	p6k.Disconnect()
	cancel()
}
