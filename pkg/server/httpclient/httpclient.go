package httpclient

import (
	"context"
)

type server interface {
	FilePlaying(ctx context.Context, file, playerIP, playerPort, playerDeviceName string) (uuid string, channels uint16, rate uint32, err error)

	PlayerReceiveStart(ctx context.Context, playerIP, playerPort string, uuid *string) (string, error)
	PlayerReceiveStop(ctx context.Context, playerIP, playerPort string) error
	PlayerPlay(ctx context.Context, playerIP, uuid, playerDeviceName string, channels, rate uint32) (err error)
	PlayerPause(ctx context.Context, playerIP, playerDeviceName string) (err error)
	PlayerStop(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid string) (err error)

	StartFileRecoding(ctx context.Context, recorderIP, recorderDeviceName string, channels, rate uint32, receivePort, file string) (err error)
	StopFileRecoding(ctx context.Context, recorderIP, recorderDeviceName, receivePort string) error

	RecorderStart(ctx context.Context, recorderIP, recorderDeviceName string, channels, rate uint32, dstAddr string) error
	RecoderStop(ctx context.Context, recorderIP, recorderDeviceName string) error

	PlayFromRecorder(ctx context.Context, playerIP, playerPort, playerDeviceName string, channels, rate uint32, recorderIP, recorderDeviceName string) (uuid string, err error)
	StopFromRecorder(ctx context.Context, playerIP, playerPort, playerDeviceName, uuid, recorderIP, recorderDeviceName string) error
}
