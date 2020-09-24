package server

import (
	"context"
	"io"
	"net"
)

// Server struct for send data over UDP connection
type Server struct {
	buffSize int
}

// TurnOn udp server
func (s *Server) TurnOn(dstAddr string) (connection io.ReadWriteCloser, err error) {
	if destinationAddress, err := net.ResolveUDPAddr("udp", dstAddr); err == nil {
		connection, err = net.DialUDP("udp", nil, destinationAddress)
	}
	return
}

// Send start sendinging data over port
func (s *Server) Send(ctx context.Context, dstAddr string, r io.Reader) (err error) {
	connection, err := s.TurnOn(dstAddr)
	if err != nil {
		return
	}

	go func() {
		outputBytes := make([]byte, s.buffSize)
		defer func() {
			connection.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				l, err := r.Read(outputBytes)
				if err != nil {
					return
				}
				connection.Write(outputBytes[:l])
			}
		}
	}()
	return
}

// NewServer return UDP server
func NewServer(buffSize int) *Server {
	return &Server{
		buffSize: buffSize,
	}
}
