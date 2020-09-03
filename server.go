package notify

import (
	"context"
	"io"
)

type Server interface {
	Connect(ClientInfo, Transport) (*Client, error)
	DisconnectClient(client *Client)
	Handle(client *Client, r io.Reader) error
	Run(ctx context.Context)
}
