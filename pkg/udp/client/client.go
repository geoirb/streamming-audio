package server

import (
	"net"
)

// ClientUDP struct for receiving data over UDP connection
type ClientUDP struct {
	connection *net.UDPConn
	localAddr  string

	c chan []byte
}

// Connect to UDP server
func (s *ClientUDP) Connect() (err error) {
	var localAddress *net.UDPAddr
	localAddress, err = net.ResolveUDPAddr("udp", s.localAddr)
	s.connection, err = net.DialUDP("udp", nil, localAddress)
	return
}

// StartReceive start receiving data over UDP connection
func (s *ClientUDP) StartReceive() {
	var data []byte
	go func() {
		length, _, _ := s.connection.ReadFromUDP(data)
		s.c <- data[:length]
	}()
}

// GetDataChan return chan for receiving data
func (s *ClientUDP) GetChan() (c <-chan []byte) {
	return s.c
}

// Disconnect udp connection
func (s *ClientUDP) Disconnect() {
	s.connection.Close()
	close(s.c)
}

// NewClientUDP return UDP client
func NewClientUDP(localAddr string) *ClientUDP {
	return &ClientUDP{
		localAddr: localAddr,
		c:         make(chan []byte, 1),
	}
}
