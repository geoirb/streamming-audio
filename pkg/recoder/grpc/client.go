package grpc

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

// StartRecode rpc request for start recode and send audio signal on server
func (c *Client) StartRecode(ctx context.Context, destAddr, recoderIP, deviceName string, channels, rate int) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, recoderIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}

	_, err = NewRecoderClient(conn).
		StartRecode(
			ctx,
			&StartRecodeRequest{
				DestAddr:   destAddr,
				DeviceName: deviceName,
				Channels:   uint32(channels),
				Rate:       uint32(rate),
			})
	if err != nil {
		return
	}
	return
}

// StopRecode rpc request for stop record and send audio signal
func (c *Client) StopRecode(ctx context.Context, destAddr, recoderIP string) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, recoderIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}

	_, err = NewRecoderClient(conn).
		StopRecode(
			ctx,
			&StopRecodeRequest{
				DestAddr: destAddr,
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
