package server

import (
	"context"
	"net"
)

// ClientUDP struct for receiving data over UDP connection
type ClientUDP struct {
	connection *net.UDPConn
	localAddr  string
	buffSize   int

	c chan []byte
}

// Connect to UDP server
func (s *ClientUDP) Connect() (err error) {
	localAddress, _ := net.ResolveUDPAddr("udp", s.localAddr)
	s.connection, _ = net.ListenUDP("udp", localAddress)
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
func NewClientUDP(localAddr string, buffSize int) *ClientUDP {
	return &ClientUDP{
		localAddr: localAddr,
		buffSize:  buffSize,
		c:         make(chan []byte, 1),
	}
}
