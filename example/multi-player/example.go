package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/geoirb/audio-service/pkg/server/httpclient"
)

type playerInfo struct {
	Start  bool
	UUID   string
	IP     string
	Port   string
	Device string
	File   string
}

var player map[string]playerInfo = map[string]playerInfo{
	"1": playerInfo{
		IP:     "127.0.0.1",
		Port:   "8081",
		Device: "hw:1,0",
		File:   "/home/geo/go/src/github.com/geoirb/audio-service/example/multi-player/test.wav",
	},
	"2": playerInfo{
		IP:     "127.0.0.1",
		Port:   "8081",
		Device: "hw:0,0",
		File:   "/home/geo/go/src/github.com/geoirb/audio-service/example/multi-player/test.wav",
	},
	"3": playerInfo{
		IP:     "127.0.0.1",
		Port:   "8081",
		Device: "hw:0,0",
		File:   "/home/geo/go/src/github.com/geoirb/audio-service/example/multi-player/test.wav",
	},
}

func main() {
	cli := httpclient.NewClient("localhost:8000")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Input number of player")
	for scanner.Scan() {
		num := scanner.Text()
		p, isExist := player[num]
		if !isExist {
			fmt.Printf("player num %v not exist\n", num)
		}
		if !p.Start {
			if uuid, _, _, err := cli.FilePlay(context.Background(), p.File, p.IP, p.Port, p.Device); err == nil {
				p.UUID = uuid
				p.Start = true
				player[num] = p
				continue
			} else {
				fmt.Println(err)
			}
		} else {
			if err := cli.FileStop(context.Background(), p.IP, p.Port, p.Device, p.UUID); err == nil {
				p.Start = false
				player[num] = p
				continue
			} else {
				fmt.Println(err)
			}
		}
	}
}
