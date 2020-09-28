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

// StartRecord rpc request for start record and send audio signal on server
func (c *Client) StartRecord(ctx context.Context, destAddr, recorderIP, deviceName string, channels, rate uint32) (err error) {
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
		StartRecord(
			ctx,
			&StartRecordRequest{
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

// StopRecord rpc request for stop record and send audio signal
func (c *Client) StopRecord(ctx context.Context, recorderIP, deviceName string) (err error) {
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
		StopRecord(
			ctx,
			&StopRecordRequest{
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
