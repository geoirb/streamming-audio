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

// State todo
func (c *Client) State(ctx context.Context, ip string) (devices []string, err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, ip, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	if err != nil {
		return
	}
	defer conn.Close()

	if res, err := NewRecorderClient(conn).
		State(
			ctx,
			&StateRequest{},
		); err == nil {
		devices = res.Devices
	}
	return
}

// Start rpc request for start record and send audio signal on server
func (c *Client) Start(ctx context.Context, destAddr, recorderIP, deviceName string, channels, rate uint32) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, recorderIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = NewRecorderClient(conn).
		Start(
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

// Stop rpc request for stop record and send audio signal
func (c *Client) Stop(ctx context.Context, recorderIP, deviceName string) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, recorderIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = NewRecorderClient(conn).
		Stop(
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
