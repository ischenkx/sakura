package notify

import (
	"io"
	"time"
)

type ServerConfig struct {
	Broker        Broker
	CleanInterval time.Duration
	DataHandler   func(client *Client, r io.Reader) error
}