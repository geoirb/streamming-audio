package player

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Client rpc controller
type Client struct {
	hostLayout  string
	controlPort string
}

// StartReceive rpc request for start receive signal from server
func (c *Client) StartReceive(ctx context.Context, playerIP, receivePort string, UUID *string) (storageUUID string, err error) {
	conn, err := grpc.Dial(
		fmt.Sprintf(c.hostLayout, playerIP, c.controlPort),
		// todo
		grpc.WithInsecure(),
	)
	defer conn.Close()
	if err != nil {
		return
	}
	req := &StartReceiveRequest{
		Port: receivePort,
	}
	if UUID != nil {
		req.StorageUUID = &wrapperspb.StringValue{
			Value: *UUID,
		}
	}

	res, err := NewPlayerClient(conn).
		StartReceive(
			ctx,
			req,
		)
	if err == nil {
		storageUUID = res.StorageUUID
	}
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
func (c *Client) StartPlay(ctx context.Context, playerIP, deviceName, storageUUID string, channels, rate uint32) (err error) {
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
				DeviceName:  deviceName,
				Channels:    uint32(channels),
				Rate:        rate,
				StorageUUID: storageUUID,
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

// ClearStorage rpc request for clear audio storage pn player
func (c *Client) ClearStorage(ctx context.Context, playerIP, storageUUID string) (err error) {
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
		ClearStorage(
			ctx,
			&ClearStorageRequest{
				StorageUUID: storageUUID,
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
