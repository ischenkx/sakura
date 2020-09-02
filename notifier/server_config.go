package notifier

import (
	"time"
)

type Config struct {
	Broker        Broker
	CleanInterval time.Duration
	DataHandler   DataHandler
}