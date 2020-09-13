package server

import (
	"context"
	"fmt"
	"io"
	"net"
)

// Server struct for send data over UDP connection
type Server struct {
	buffSize int
}

// Send start sendinging data over port
func (s *Server) Send(ctx context.Context, dstAddr string, r io.Reader) (err error) {
	var (
		destinationAddress *net.UDPAddr
		connection         *net.UDPConn
	)

	if destinationAddress, err = net.ResolveUDPAddr("udp", dstAddr); err != nil {
		return
	}
	if connection, err = net.DialUDP("udp", nil, destinationAddress); err != nil {
		return
	}

	go func() {
		defer func() {
			connection.Close()
		}()

		outputBytes := make([]byte, s.buffSize)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				l, err := r.Read(outputBytes)
				if err != nil {
					return
				}
				fmt.Println(outputBytes)
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
