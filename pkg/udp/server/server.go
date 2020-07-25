package server

import (
	"net"
)

// ServerUDP struct for send data over UDP connection
type ServerUDP struct {
	connection *net.UDPConn
	dstAddr    string
}

// Start configuring and starting UDP server
func (s *ServerUDP) Start() (err error) {
	var destinationAddress *net.UDPAddr
	destinationAddress, err = net.ResolveUDPAddr("udp", s.dstAddr)
	s.connection, err = net.DialUDP("udp", nil, destinationAddress)
	return
}

// Send data on UDP connection
func (s *ServerUDP) Send(data []byte) (err error) {
	_, err = s.connection.Write(data)
	return
}

// Shutdown UDP server
func (s *ServerUDP) Shutdown() {
	s.connection.Close()
}

// NewServerUDP return UDP server
func NewServerUDP(dstAddr string) *ServerUDP {
	return &ServerUDP{
		dstAddr: dstAddr,
	}
}
