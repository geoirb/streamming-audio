package server

import (
	"net"
)

// Server struct for send data over UDP connection
type Server struct {
	connection *net.UDPConn
}

// Send data on UDP connection
func (s *Server) Send(data []byte) (err error) {
	_, err = s.connection.Write(data)
	return
}

// Shutdown UDP server
func (s *Server) Shutdown() error {
	return s.connection.Close()
}

// NewServer return UDP server
func NewServer(dstAddr string) (s *Server, err error) {
	if destinationAddress, err := net.ResolveUDPAddr("udp", dstAddr); err == nil {
		s.connection, err = net.DialUDP("udp", nil, destinationAddress)
	}
	return
}
