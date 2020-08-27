package server

import (
	"context"
	"net"
)

// ClientUDP struct for receiving data over UDP connection
type ClientUDP struct {
	connection *net.UDPConn
	port       string
	buffSize   int

	c chan []byte
}

// Connect to UDP server
func (s *ClientUDP) Connect() (err error) {
	var clientAddress *net.UDPAddr
	if clientAddress, err = net.ResolveUDPAddr("udp", s.port); err == nil {
		s.connection, err = net.ListenUDP("udp", clientAddress)
	}
	return
}

// StartReceive start receiving data over UDP connection
func (s *ClientUDP) StartReceive(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			inputBytes := make([]byte, s.buffSize)
			l, _, err := s.connection.ReadFromUDP(inputBytes)
			if err != nil {
				return
			}
			s.c <- inputBytes[:l]
		}
	}
}

// Data return data chan
func (s *ClientUDP) Data() <-chan []byte {
	return s.c
}

// Disconnect udp connection
func (s *ClientUDP) Disconnect() {
	s.connection.Close()
	close(s.c)
}

// NewClientUDP return UDP client
func NewClientUDP(port string, buffSize int) *ClientUDP {
	return &ClientUDP{
		port:     port,
		buffSize: buffSize,
		c:        make(chan []byte, 1),
	}
}
