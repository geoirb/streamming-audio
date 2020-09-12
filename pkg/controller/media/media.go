package media

import (
	context "context"
	fmt "fmt"
	"sync"

	grpc "google.golang.org/grpc"
)

// Controller rpc controller
type Controller struct {
	hostLayout string // "%s:%d"
	port       string

	mutex      sync.Mutex
	connection map[string]*grpc.ClientConn
}

// StartReceive rpc request for start receive and play audio signal
func (c *Controller) StartReceive(ctx context.Context, ip, port, deviceName string, channels, rate uint32) (err error) {
	host := fmt.Sprintf(c.hostLayout, ip, c.port)
	conn, err := grpc.Dial(
		host,
		grpc.WithInsecure(),
	)
	if err != nil {
		return
	}

	client := NewMediaClient(conn)
	_, err = client.StartReceive(
		ctx,
		&StartReceiveRequest{
			Port:       port,
			DeviceName: deviceName,
			Channels:   uint32(channels),
			Rate:       rate,
		})
	if err != nil {
		return
	}

	c.mutex.Lock()
	c.connection[host] = conn
	c.mutex.Unlock()
	return
}

// StopReceive rpc request for stop receive and play audio signal
func (c *Controller) StopReceive(ctx context.Context, ip, port string) (err error) {
	host := fmt.Sprintf(c.hostLayout, ip, c.port)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	conn, isExist := c.connection[host]
	if !isExist {
		err = fmt.Errorf("client %s not exist", ip)
		return
	}
	client := NewMediaClient(conn)
	_, err = client.StopReceive(
		ctx,
		&StopReceiveRequest{
			Port: port,
		},
	)
	if err != nil {
		return
	}

	conn.Close()
	delete(c.connection, host)
	return
}

// NewMediaController ...
func NewMediaController(hostLayout, port string) *Controller {
	return &Controller{
		hostLayout: hostLayout,
		port:       port,
		connection: make(map[string]*grpc.ClientConn),
	}
}
