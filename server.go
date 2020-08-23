package notify

import "io"

type Server interface {
	Connect(app string, info ClientInfo) (*Client, error)
	DisconnectClient(client *Client)
	Handle(client *Client, r io.Reader) error
}
