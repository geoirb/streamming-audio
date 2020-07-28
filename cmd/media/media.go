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
	localAddr  = ":8080"
	deviceName = "default"
	channels   = 1
	rate       = 4100
	sizeData   = 10
	buffSize   = 1024
)

func main() {
	udpClt := udp.NewClientUDP(localAddr, buffSize)
	udpClt.Connect()
	defer udpClt.Disconnect()

	d4c := device.NewDevice(
		deviceName,
		channels,
		rate,
	)
	d4c.Connect()
	defer d4c.Disconnect()

	c7r := converter.NewConverter()

	m := media.NewMedia(
		udpClt,
		d4c,
		c7r,
		sizeData,
	)
	ctx, cancel := context.WithCancel(context.Background())
	go m.Repicenting(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	sig := <-c
	fmt.Printf("received signal, exiting signal %v\n", sig)
	cancel()
}
