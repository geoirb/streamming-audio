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

// StartRecode rpc request for start recode and send audio signal on server
func (c *Client) StartRecode(ctx context.Context, ip, port, deviceName string, channels, rate uint32) (err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, ip, c.port),
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

	_, err = NewRecoderClient(conn).
		StopRecode(
			ctx,
			&StopRecodeRequest{
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
