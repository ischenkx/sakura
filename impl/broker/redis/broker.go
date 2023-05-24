package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"sakura/common/data/codec"
	"sakura/core/broker"
)

type Broker[T any] struct {
	client redis.UniversalClient
	codec  codec.Binary[T]
}

func (b *Broker[T]) Push(ctx context.Context, channel string, message T) error {
	payload, err := b.codec.Encoder().Convert(message)
	if err != nil {
		return err
	}
	return b.client.Publish(ctx, channel, payload).Err()
}

func (b *Broker[T]) PubSub(ctx context.Context) broker.PubSub[T] {
	pubsub := b.client.Subscribe(ctx)
	return &PubSub[T]{
		pubsub: pubsub,
		codec:  b.codec,
	}
}

type PubSub[T any] struct {
	pubsub *redis.PubSub
	codec  codec.Binary[T]
}

func (p *PubSub[T]) Subscribe(ctx context.Context, channels ...string) error {
	if len(channels) == 0 {
		return nil
	}
	return p.pubsub.Subscribe(ctx, channels...)
}

func (p *PubSub[T]) Unsubscribe(ctx context.Context, channels ...string) error {
	if len(channels) == 0 {
		return nil
	}
	return p.pubsub.Unsubscribe(ctx, channels...)
}

func (p *PubSub[T]) Channel(ctx context.Context) (<-chan broker.Message[T], error) {
	messages := p.pubsub.Channel()
	outputs := make(chan broker.Message[T], 512)

	go func(ctx context.Context, ch chan<- broker.Message[T]) {
		<-ctx.Done()
		close(ch)
	}(ctx, outputs)

	go func(ctx context.Context, from <-chan *redis.Message, to chan<- broker.Message[T]) {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop

			case rawMessage := <-messages:
				payload := []byte(rawMessage.Payload)
				message, err := p.codec.Decoder().Convert(payload)
				if err != nil {
					log.Println("failed to decode the message:", err)
					continue
				}

				to <- broker.Message[T]{
					Channel: rawMessage.Channel,
					Data:    []byte(message),
				}
			}
		}
	}(ctx, messages, outputs)

	return outputs, nil
}

func (p *PubSub[T]) Clear(ctx context.Context) error {
	return p.pubsub.Unsubscribe(ctx)
}
