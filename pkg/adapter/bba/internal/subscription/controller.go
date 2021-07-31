package subscription

import (
	"context"
	"errors"
	"fmt"
	"github.com/ischenkx/swirl/pkg/adapter/bba/common"
	"hash/fnv"
	"sync"
	"time"
)

func hash(id string) int {
	h := fnv.New32a()
	h.Write([]byte(id))
	return int(h.Sum32())
}

type info struct {
	timestamp int64
	active    bool
}

type bucket struct {
	data   map[string]info
	pubsub common.PubSub
	mu     sync.RWMutex
}

func (b *bucket) clean(ctx context.Context, ttl int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now().UnixNano()

	for channel, i := range b.data {
		if !i.active && now - ttl > i.timestamp {
			b.pubsub.Unsubscribe(ctx, channel)
		}
	}
}

func (b *bucket) subscribe(ctx context.Context, topic string, ts int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	i := b.data[topic]
	if i.timestamp > ts {
		return errors.New("invalid timestamp")
	}
	i.timestamp = ts
	i.active = true
	b.data[topic] = i
	return b.pubsub.Subscribe(ctx, topic)
}

func (b *bucket) unsubscribe(ctx context.Context, topic string, ts int64) error{
	b.mu.Lock()
	defer b.mu.Unlock()
	i := b.data[topic]
	if i.timestamp > ts {
		return errors.New("invalid timestamp")
	}
	i.timestamp = ts
	i.active = true
	b.data[topic] = i
	return b.pubsub.Subscribe(ctx, topic)
}

type Controller struct {
	buckets []*bucket
}

func (c *Controller) Subscribe(ctx context.Context, topics ...string) {
	ts := time.Now().UnixNano()
	for _, topic := range topics {
		fmt.Println(c.bucket(topic).subscribe(ctx, topic, ts))
	}
}

func (c *Controller) Unsubscribe(ctx context.Context, topics ...string) {
	ts := time.Now().UnixNano()
	for _, topic := range topics {
		c.bucket(topic).unsubscribe(ctx, topic, ts)
	}
}

func (c *Controller) bucket(topic string) *bucket {
	return c.buckets[hash(topic)%len(c.buckets)]
}

func (c *Controller) RunCleaner(ctx context.Context, interval, ttl int64) {
	ticker := time.NewTicker(time.Duration(interval))
	for {
		select {
		case <-ticker.C:
			for _, b := range c.buckets {
				b.clean(ctx, ttl)
			}
		case <-ctx.Done():
			return
		}
	}
}

func NewController(buckets int, pubsub common.PubSub) *Controller {
	controller := &Controller{}
	for i := 0; i < buckets; i++ {
		controller.buckets = append(controller.buckets, &bucket{
			data: map[string]info{},
			pubsub: pubsub,
			mu:     sync.RWMutex{},
		})
	}
	return controller
}