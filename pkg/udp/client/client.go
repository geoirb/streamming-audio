package server

import (
	"net"
)

// ClientUDP struct for receiving data over UDP connection
type ClientUDP struct {
	connection *net.UDPConn
	localAddr  string
}

// Connect to UDP server
func (s *ClientUDP) Connect() (err error) {
	var localAddress *net.UDPAddr
	localAddress, err = net.ResolveUDPAddr("udp", s.localAddr)
	s.connection, err = net.DialUDP("udp", nil, localAddress)
	return
}

// Receive data over UDP connection
func (s *ClientUDP) Receive() (data []byte, err error) {
	var length int
	length, _, err = s.connection.ReadFromUDP(data)
	data = data[:length]
	return
}

// Disconnect udp connection
func (s *ClientUDP) Disconnect() {
	s.connection.Close()
}

// NewClientUDP return UDP client
func NewClientUDP(localAddr string) *ClientUDP {
	return &ClientUDP{
		localAddr: localAddr,
	}
}
