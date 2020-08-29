package notify

import (
	"context"
	"io"
)

type Server interface {
	//first argument - data to be passed to auth or ClientInfo
	//second - transport
	Connect(interface{}, Transport) (*Client, error)
	DisconnectClient(client *Client)
	Handle(client *Client, r io.Reader) error
	Run(ctx context.Context)
}
