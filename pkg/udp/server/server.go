package server

import (
	"net"
)

// ConnectionUDP struct for send data over UDP connection
type ConnectionUDP struct {
	connection *net.UDPConn
	dstAddr    string
}

// TurnOn configuring and starting UDP server
func (s *ConnectionUDP) TurnOn() (err error) {
	var destinationAddress *net.UDPAddr
	if destinationAddress, err = net.ResolveUDPAddr("udp", s.dstAddr); err == nil {
		s.connection, err = net.DialUDP("udp", nil, destinationAddress)
	}
	return
}

// Send data on UDP connection
func (s *ConnectionUDP) Send(data []byte) (err error) {
	_, err = s.connection.Write(data)
	return
}

// Shutdown UDP server
func (s *ConnectionUDP) Shutdown() error {
	return s.connection.Close()
}

// NewServerUDP return UDP server
func NewServerUDP(dstAddr string) *ConnectionUDP {
	return &ConnectionUDP{
		dstAddr: dstAddr,
	}
}
