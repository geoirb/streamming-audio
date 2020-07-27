package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/GeoIrb/sound-ethernet-streaming/pkg/converter"
	"github.com/GeoIrb/sound-ethernet-streaming/pkg/device"
	"github.com/GeoIrb/sound-ethernet-streaming/pkg/media"
	udp "github.com/GeoIrb/sound-ethernet-streaming/pkg/udp/client"
)

const (
	localAddr  = ":1235"
	deviceName = "default"
	channels   = 1
	rate       = 4100
	sizeData   = 10
)

func main() {
	udpClt := udp.NewClientUDP(localAddr)
	udpClt.Connect()
	defer udpClt.Disconnect()

	dvc := device.NewDevice(
		deviceName,
		channels,
		rate,
	)

	cnv := converter.NewConverter()

	m := media.NewMedia(
		udpClt,
		dvc,
		cnv,
		sizeData,
	)
	ctx, cancel := context.WithCancel(context.Background())
	go m.Repicenting(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	fmt.Printf("received signal, exiting signal %v", sig)
	cancel()
}
