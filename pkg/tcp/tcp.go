package tcp

import (
	"context"
	"fmt"
	"io"
	"net"
)

// TCP receive and send
type TCP struct {
	buffSize int
}

// TurnOnSender tcp sender
func (u *TCP) TurnOnSender(dstAddr string) (connection io.WriteCloser, err error) {
	connection, err = net.Dial("tcp", dstAddr)
	return
}

// Send start sendinging data over port
func (u *TCP) Send(ctx context.Context, dstAddr string, r io.Reader) (err error) {
	connection, err := u.TurnOnSender(dstAddr)
	if err != nil {
		return
	}

	go func() {
		outputBytes := make([]byte, u.buffSize)
		defer func() {
			connection.Close()
		}()
		i := 0
		for {
			i++
			select {
			case <-ctx.Done():
				return
			default:
				l, err := r.Read(outputBytes)
				if err != nil {
					return
				}
				fmt.Println(i)
				connection.Write(outputBytes[:l])
			}
		}
	}()
	return
}

// Receive start receiving data over port
func (u *TCP) Receive(ctx context.Context, receivePort string, w io.Writer) (err error) {
	ln, err := net.Listen("tcp", ":"+receivePort)
	if err != nil {
		return
	}

	go func() {
		connection, _ := ln.Accept()

		go func() {
			<-ctx.Done()
			connection.Close()
		}()

		i := 0
		for {
			i++
			inputBytes := make([]byte, u.buffSize)
			l, err := connection.Read(inputBytes)
			if err != nil {
				return
			}
			fmt.Println(i)
			w.Write(inputBytes[:l])
		}
	}()
	return
}

// NewTCP ...
func NewTCP(buffSize int) *TCP {
	return &TCP{
		buffSize: buffSize,
	}
}
