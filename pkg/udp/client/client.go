package server

import (
	"context"
	"net"

	"github.com/geoirb/sound-ethernet-streaming/pkg/cash"
)

// Client struct for receiving data over UDP connection
type Client struct {
	buffSize int
}

// Receive start receiving data over port
func (c *Client) Receive(ctx context.Context, port string, ca *cash.Cash) (err error) {
	var (
		clientAddress *net.UDPAddr
		connection    *net.UDPConn
	)

	if clientAddress, err = net.ResolveUDPAddr("udp", port); err != nil {
		return
	}
	if connection, err = net.ListenUDP("udp", clientAddress); err != nil {
		return
	}

	go func() {
		defer func() {
			connection.Close()
		}()

		inputBytes := make([]byte, c.buffSize)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				l, _, err := connection.ReadFromUDP(inputBytes)
				if err != nil {
					return
				}
				ca.Push(inputBytes[:l])
			}
		}
	}()

	return
}

// NewClient return UDP client
func NewClient(buffSize int) *Client {
	return &Client{
		buffSize: buffSize,
	}
}
