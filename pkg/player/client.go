package player

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

// Client rpc controller
type Client struct {
	hostLayout  string
	controlPort string
}

// StartReceive rpc request for start receive signal from server
func (c *Client) StartReceive(ctx context.Context, playerIP, receivePort string) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, playerIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}

	_, err = NewPlayerClient(conn).
		StartReceive(
			ctx,
			&StartReceiveRequest{
				Port: receivePort,
			})
	return
}

// StopReceive rpc request for stop receive signal from server
func (c *Client) StopReceive(ctx context.Context, playerIP, receivePort string) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, playerIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}

	_, err = NewPlayerClient(conn).
		StopReceive(
			ctx,
			&StopReceiveRequest{
				Port: receivePort,
			})
	return
}

// StartPlay rpc request for play audio signal
func (c *Client) StartPlay(ctx context.Context, playerIP, receivePort, deviceName string, channels, rate uint32) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, playerIP, c.controlPort),
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
				Port:       receivePort,
				DeviceName: deviceName,
				Channels:   uint32(channels),
				Rate:       rate,
			})
	return
}

// StopPlay rpc request for play audio signal
func (c *Client) StopPlay(ctx context.Context, playerIP, deviceName string) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, playerIP, c.controlPort),
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
				DeviceName: deviceName,
			},
		)
	return
}

// NewClient ...
func NewClient(hostLayout, controlPort string) *Client {
	return &Client{
		hostLayout:  hostLayout,
		controlPort: controlPort,
	}
}
