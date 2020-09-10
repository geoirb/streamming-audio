package server

import (
	"context"
	"net"
)

// Client struct for receiving data over UDP connection
type Client struct {
	buffSize int
}

// Receive start receiving data over port
func (c *Client) Receive(ctx context.Context, port string) (<-chan []byte, error) {
	var (
		data          chan []byte
		clientAddress *net.UDPAddr
		connection    *net.UDPConn
		err           error
	)
	if clientAddress, err = net.ResolveUDPAddr("udp", port); err != nil {
		return data, err
	}
	if connection, err = net.ListenUDP("udp", clientAddress); err != nil {
		return data, err
	}

	data = make(chan []byte, c.buffSize)

	go func() {
		defer func() {
			close(data)
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
				data <- inputBytes[:l]
			}
		}
	}()

	return data, err
}

// NewClient return UDP client
func NewClient(port string, buffSize int) *Client {
	return &Client{
		buffSize: buffSize,
	}
}
