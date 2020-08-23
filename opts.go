package notify

import "io"

type JoinOptions struct {
	Users []string
	Clients []string
	Channels []string
}

type LeaveOptions struct {
	Users []string
	Clients []string
	Channels []string
	All bool
}

type SendOptions struct {
	Users []string
	Clients []string
	Channels []string
	Message
}

type IncomingDataOptions struct {
	Client *Client
	Reader io.Reader
}

type ClientConnectionOptions struct {
	Client *Client
}

type ClientDisconnectionOptions struct {
	ID string
}
