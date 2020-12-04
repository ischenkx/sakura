package transport

import "io"

type Transport interface {
	io.Writer
	io.Closer
	State() State
}
