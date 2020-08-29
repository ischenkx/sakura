package notify

import (
	"github.com/RomanIschenko/notify/broker"
	"time"
)

type ServerConfig struct {
	Auth          Auth
	Broker        broker.Broker
	CleanInterval time.Duration
	DataHandler	  DataHandler
}