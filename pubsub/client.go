package pubsub

import "sync"

type Client interface {
	ID() string
	User() string
	Data() *sync.Map
}
