package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"sakura/common/data/codec"
	"sakura/core/broker"
)

type Broker[Message any] struct {
	client redis.UniversalClient
	codec  codec.Binary[Message]
}

func (b *Broker[Message]) Push(ctx context.Context, channel string, message Message) error {
	payload, err := b.codec.Encoder().Convert(message)
	if err != nil {
		return err
	}
	return b.client.Publish(ctx, channel, payload).Err()
}

func (b *Broker[Message]) PubSub(ctx context.Context) broker.PubSub[Message] {
	pubsub := b.client.Subscribe(ctx)
	return &PubSub[Message]{
		pubsub: pubsub,
		codec:  b.codec,
	}
}

type PubSub[Message any] struct {
	pubsub *redis.PubSub
	codec  codec.Binary[Message]
}

func (p *PubSub[Message]) Subscribe(ctx context.Context, channels ...string) error {
	if len(channels) == 0 {
		return nil
	}
	return p.pubsub.Subscribe(ctx, channels...)
}

func (p *PubSub[Message]) Unsubscribe(ctx context.Context, channels ...string) error {
	if len(channels) == 0 {
		return nil
	}
	return p.pubsub.Unsubscribe(ctx, channels...)
}

func (p *PubSub[Message]) Channel(ctx context.Context) (<-chan Message, error) {
	messages := p.pubsub.Channel()
	outputs := make(chan Message, 512)

	go func(ctx context.Context, ch chan<- Message) {
		<-ctx.Done()
		close(ch)
	}(ctx, outputs)

	go func(ctx context.Context, from <-chan *redis.Message, to chan<- Message) {
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

				to <- message
			}
		}
	}(ctx, messages, outputs)

	return outputs, nil
}

func (p *PubSub[Message]) Clear(ctx context.Context) error {
	return p.pubsub.Unsubscribe(ctx)
}
