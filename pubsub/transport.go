package pubsub

import "io"

type TransportState int

const (
	OpenTransport TransportState = iota
	ClosedTransport
)

type Transport interface {
	io.Writer
	io.Closer
	State() TransportState
}
