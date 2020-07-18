package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

func main() {
	conn, _ := net.Dial("udp", ":10001")
	defer conn.Close()

	data, err := ioutil.ReadFile("/home/geo/go/src/github.com/GeoIrb/sound-ethernet-streaming/audio/Anacondaz.mp3")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	for {
		fmt.Println("File reading", len(data))
		conn.Write(data[1024:2049])
		time.Sleep(time.Second * 1)
	}
}
