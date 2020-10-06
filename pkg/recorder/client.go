package recorder

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

// StartSend rpc request for start record and send audio signal on server
func (c *Client) StartSend(ctx context.Context, destAddr, recorderIP, deviceName string, channels, rate uint32) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, recorderIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}

	_, err = NewRecorderClient(conn).
		StartSend(
			ctx,
			&StartSendRequest{
				DeviceName: deviceName,
				Channels:   channels,
				Rate:       rate,
				DestAddr:   destAddr,
			})
	if err != nil {
		return
	}
	return
}

// StopSend rpc request for stop record and send audio signal
func (c *Client) StopSend(ctx context.Context, recorderIP, deviceName string) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, recorderIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}

	_, err = NewRecorderClient(conn).
		StopSend(
			ctx,
			&StopSendRequest{
				DeviceName: deviceName,
			},
		)
	if err != nil {
		return
	}
	return
}

// NewClient ...
func NewClient(hostLayout, controlPort string) *Client {
	return &Client{
		hostLayout:  hostLayout,
		controlPort: controlPort,
	}
}
