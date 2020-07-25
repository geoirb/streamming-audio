package server

import (
	"net"
)

type ClientUDP struct {
	connection *net.UDPConn
	localAddr  string
}

func (s *ClientUDP) Connect() (err error) {
	var localAddress *net.UDPAddr
	localAddress, err = net.ResolveUDPAddr("udp", s.localAddr)
	s.connection, err = net.DialUDP("udp", nil, localAddress)
	return
}

func (s *ClientUDP) Read() (data []byte, err error) {
	var length int
	length, _, err = s.connection.ReadFromUDP(data)
	data = data[:length]
	return
}

func (s *ClientUDP) Disconnect() {
	s.connection.Close()
}
