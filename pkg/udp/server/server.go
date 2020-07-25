package server

import (
	"net"
)

type ServerUDP struct {
	connection *net.UDPConn
	dstAddr    string
}

func (s *ServerUDP) Connect() (err error) {
	var destinationAddress *net.UDPAddr
	destinationAddress, err = net.ResolveUDPAddr("udp", s.dstAddr)
	s.connection, err = net.DialUDP("udp", nil, destinationAddress)
	return
}

func (s *ServerUDP) Write(data []byte) (err error) {
	_, err = s.connection.Write(data)
	return
}

func (s *ServerUDP) Disconnect() {
	s.connection.Close()
}
