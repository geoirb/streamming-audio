package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr, _ := net.ResolveUDPAddr("udp", ":10001")
	serverConn, _ := net.ListenUDP("udp", serverAddr)
	defer serverConn.Close()

	buf := make([]byte, 1024)

	for {
		n, _, _ := serverConn.ReadFromUDP(buf)
		fmt.Println(buf[0:n])
	}
}
