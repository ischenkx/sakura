package redibroker

import (
	"bytes"
	"context"
	"github.com/RomanIschenko/notify/cluster/broker"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var logger = logrus.WithField("source", "redis")

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

type Broker struct {
	client redis.UniversalClient
	ps *redis.PubSub
	manager *channelManager
	brokerId string
	handlers map[string]func(e broker.Event)
	mu sync.RWMutex
}

func (p *Broker) Publish(channels []string, e broker.Event) error {
	if len(channels) == 0 {
		return nil
	}
	e.BrokerID = p.brokerId
	pipeline := p.client.Pipeline()
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buffer.Reset()
		bufferPool.Put(buffer)
	}()

	err := encodeEvent(e, buffer)
	if err != nil {
		return err
	}
	for _, ch := range channels {
		pipeline.Publish(context.Background(), ch, buffer.Bytes())
	}
	_, err = pipeline.Exec(context.Background())
	return err
}

func (p *Broker) Subscribe(chs []string, t int64) {
	addedChannels := p.manager.addChannels(chs, t)
	if err := p.ps.Subscribe(context.Background(), addedChannels...); err != nil {
		logger.Errorf("redis subscription error: %s", err)
	}
}

func (p *Broker) ID() string {
	return p.brokerId
}


func (p *Broker) Unsubscribe(chs []string, t int64) {
	p.manager.delChannels(chs, t)
}

func (p *Broker) Handle(h func(broker.Event)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	uid := uuid.New().String()
	p.handlers[uid] = h
}

func (p *Broker) readEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.ps.Close()
			return
		case mes := <-p.ps.Channel():
			e, err := decodeEvent([]byte(mes.Payload))
			if err == nil {
				e.Source = mes.Channel
				if e.BrokerID == p.brokerId {
					continue
				}
				p.mu.RLock()
				for _, h := range p.handlers {
					h(e)
				}
				p.mu.RUnlock()
			} else {
				logger.Infof("decoding event error: %s", err)
			}
		}

	}
}

func (p *Broker) runCleaner(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if unsubs := p.manager.clean(); len(unsubs) > 0 {
				if err := p.ps.Unsubscribe(ctx, unsubs...); err != nil {
					logger.Debugf("redis unsubscribe error: %s", err)
				}
			}
		}
	}
}

func newBroker(ctx context.Context, client redis.UniversalClient, brokerId string) *Broker {
	ps := client.Subscribe(ctx)
	x := &Broker{
		ps:       ps,
		brokerId: brokerId,
		handlers: map[string]func(e broker.Event){},
		client:   client,
		manager:  newChannelManager(),
	}
	go x.runCleaner(ctx)
	go x.readEvents(ctx)
	return x
}