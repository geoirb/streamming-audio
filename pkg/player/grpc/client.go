package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// Client rpc controller
type Client struct {
	hostLayout string
	port       string
}

// StartPlay rpc request for start receive and play audio signal
func (c *Client) StartPlay(ctx context.Context, ip, port, deviceName string, channels, rate uint32) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, ip, c.port),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}

	_, err = NewPlayerClient(conn).
		StartPlay(
			ctx,
			&StartPlayRequest{
				Port:       port,
				DeviceName: deviceName,
				Channels:   uint32(channels),
				Rate:       rate,
			})
	if err != nil {
		return
	}
	return
}

// StopPlay rpc request for stop receive and play audio signal
func (c *Client) StopPlay(ctx context.Context, ip, port string) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, ip, c.port),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}

	_, err = NewPlayerClient(conn).
		StopPlay(
			ctx,
			&StopPlayRequest{
				Port: port,
			},
		)
	if err != nil {
		return
	}
	return
}

// NewClient ...
func NewClient(hostLayout, port string) *Client {
	return &Client{
		hostLayout: hostLayout,
		port:       port,
	}
}
