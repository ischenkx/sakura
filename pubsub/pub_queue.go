package pubsub

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub/publication"
)

const DefaultPubQueueBufferSize = 4096
const DefaultPubQueueWriters = 4
const DefaultPubQueueRPW = 2

type PubQueueConfig struct {
	BufferSize, Writers, ReadersPerWriter int
}

func (cfg PubQueueConfig) validate() PubQueueConfig {
	if cfg.BufferSize <= 0 {
		cfg.BufferSize = DefaultPubQueueBufferSize
	}
	if cfg.Writers <= 0 {
		cfg.Writers = DefaultPubQueueWriters
	}
	if cfg.ReadersPerWriter <= 0 {
		cfg.ReadersPerWriter = DefaultPubQueueRPW
	}
	return cfg
}

type clientPublication struct {
	client *Client
	pub    publication.Publication
}

type pubQueue struct {
	writers []chan clientPublication
	config PubQueueConfig
}

func (b pubQueue) Enqueue(c *Client, pub publication.Publication) {
	idx := c.Hash() % len(b.writers)
	writer := b.writers[idx]
	writer <- clientPublication{c, pub}
}

func (b pubQueue) Start(ctx context.Context) {
	for _, writer := range b.writers {
		for i := 0; i < b.config.ReadersPerWriter; i++ {
			go func(ctx context.Context, w chan clientPublication) {
				for {
					select {
					case <-ctx.Done():
						return
					case cpub := <-w:
						cpub.client.Publish(cpub.pub)
					}
				}
			}(ctx, writer)
		}
	}
}

func newPubQueue(cfg PubQueueConfig) pubQueue {
	cfg = cfg.validate()

	writers := make([]chan clientPublication, cfg.Writers)
	for i := range writers {
		writers[i] = make(chan clientPublication, cfg.BufferSize)
	}
	return pubQueue{writers, cfg}
}