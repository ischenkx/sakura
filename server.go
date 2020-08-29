package notify

import (
	"context"
	"io"
)

type Server interface {
	Connect(interface{}, Transport) (*Client, error)
	DisconnectClient(client *Client)
	Handle(client *Client, r io.Reader) error
	Run(ctx context.Context)
}
