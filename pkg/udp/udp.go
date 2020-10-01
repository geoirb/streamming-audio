package udp

import (
	"context"
	"io"
	"net"
)

// UDP receive and send
type UDP struct {
	buffSize int
}

// TurnOnSender udp sender
func (u *UDP) TurnOnSender(dstAddr string) (connection io.WriteCloser, err error) {
	var destinationAddress *net.UDPAddr
	if destinationAddress, err = net.ResolveUDPAddr("udp", dstAddr); err != nil {
		return
	}
	connection, err = net.DialUDP("udp", nil, destinationAddress)
	return
}

// Send start sendinging data over port
func (u *UDP) Send(ctx context.Context, dstAddr string, r io.Reader) (err error) {
	connection, err := u.TurnOnSender(dstAddr)
	if err != nil {
		return
	}

	go func() {
		outputBytes := make([]byte, u.buffSize)
		defer func() {
			connection.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				l, err := r.Read(outputBytes)
				if err != nil {
					return
				}
				connection.Write(outputBytes[:l])
			}
		}
	}()
	return
}

// TurnOnReceiver udp receiver
func (u *UDP) TurnOnReceiver(receivePort string) (connection io.ReadWriteCloser, err error) {
	var receiveAddress *net.UDPAddr
	if receiveAddress, err = net.ResolveUDPAddr("udp", ":"+receivePort); err != nil {
		return
	}
	connection, err = net.ListenUDP("udp", receiveAddress)
	return
}

// Receive start receiving data over port
func (u *UDP) Receive(ctx context.Context, receivePort string, w io.Writer) (err error) {
	connection, err := u.TurnOnReceiver(receivePort)
	if err != nil {
		return
	}

	go func() {
		<-ctx.Done()
		connection.Close()
	}()

	go func() {
		for {
			inputBytes := make([]byte, u.buffSize)
			l, err := connection.Read(inputBytes)
			if err != nil {
				return
			}
			w.Write(inputBytes[:l])
		}
	}()
	return
}

// NewUDP ...
func NewUDP(buffSize int) *UDP {
	return &UDP{
		buffSize: buffSize,
	}
}
