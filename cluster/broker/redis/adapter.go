package redibroker

import (
	"context"
	"github.com/RomanIschenko/notify/cluster/broker"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"sync"
)

type Adapter struct {
	redis redis.UniversalClient
	mu sync.RWMutex
	id string
}

func (b *Adapter) Broker(ctx context.Context) broker.Broker {
	return newBroker(ctx, b.redis, b.id)
}

func (b *Adapter) ID() string {
	return b.id
}

func New(c redis.UniversalClient) *Adapter {
	return &Adapter{
		redis:    c,
		id:       uuid.New().String(),
	}
}