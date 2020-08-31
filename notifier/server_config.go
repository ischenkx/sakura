package notifier

import (
	"time"
)

type Config struct {
	Auth          Auth
	Broker        Broker
	CleanInterval time.Duration
	DataHandler   DataHandler
}