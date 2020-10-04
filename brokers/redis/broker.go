package redibroker

import (
	"context"
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/brokers"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"io"
	"sync"
	"time"
)

const (
	globalChannel = "global_channel"
)

type appEvent struct {
	op string
	app string
}

type Broker struct {
	client redis.UniversalClient
	id string
	mu sync.RWMutex
	starter chan struct{}
	apps chan appEvent
	handlers map[string]func(notify.BrokerEvent)
	protocol *brokers.RBP
}

func (b *Broker) Handle(h func(e notify.BrokerEvent)) io.Closer {
	if h == nil {
		return HandlerCloser{}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	uid := uuid.New().String()
	b.handlers[uid] = h
	return HandlerCloser{b, uid}
}

func (b *Broker) Emit(e notify.BrokerEvent) {
	if b.client == nil {
		return
	}
	e.BrokerID = b.id
	e.Time = time.Now().UnixNano()
	marshalledMessage, err := b.protocol.Encode(e)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch e.Event {
	case notify.PublishEvent, notify.SubscribeEvent, notify.UnsubscribeEvent, notify.ConnectEvent, notify.DisconnectEvent:
		b.client.Publish(context.Background(), e.AppID, marshalledMessage)
	case notify.BrokerAppUpEvent, notify.BrokerAppDownEvent:
		b.apps <- appEvent{e.Event, e.AppID}
		b.client.Publish(context.Background(), globalChannel, marshalledMessage)
	default:
		b.client.Publish(context.Background(), globalChannel, marshalledMessage)
	}
}

func (b *Broker) Run(ctx context.Context, gc int) {
	b.starter <- struct{}{}
	defer func() {
		<-b.starter
	}()
	pubsub := b.client.Subscribe(ctx)
	defer func() {
		pubsub.Close()
	}()
	if err := pubsub.Subscribe(ctx, globalChannel); err != nil {
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(gc)
	for i := 0; i < gc; i++ {
		go func(ctx context.Context, wg *sync.WaitGroup, b *Broker, pubsub *redis.PubSub) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case mes := <-pubsub.Channel():
					e, err := b.protocol.Decode([]byte(mes.Payload))
					if err != nil {
						fmt.Println("error while reading message", err)
						continue
					}
					if e.BrokerID == b.id {
						continue
					}
					b.mu.RLock()
					for _, h := range b.handlers {
						h(e)
					}
					b.mu.RUnlock()
				case appOp := <-b.apps:
					if appOp.op == notify.BrokerAppUpEvent {
						pubsub.Subscribe(context.Background(), appOp.app)
					}
					if appOp.op == notify.BrokerAppDownEvent {
						pubsub.Unsubscribe(context.Background(), appOp.app)
					}
				}
			}
		}(ctx, wg, b, pubsub)
	}
	wg.Wait()
}

func (b *Broker) ID() string {
	return b.id
}

func New(client redis.UniversalClient) *Broker {
	return &Broker{
		client:   client,
		id:       uuid.New().String(),
		starter:  make(chan struct{}, 1),
		apps:     make(chan appEvent, 128),
		protocol: brokers.NewRBP(),
		handlers: map[string]func(notify.BrokerEvent){},
	}
}