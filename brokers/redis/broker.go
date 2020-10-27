package redibroker

import (
	"context"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/brokers"
	"github.com/RomanIschenko/notify/events"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

const TopicsCleanInterval = time.Minute * 4

var logger = logrus.WithField("source", "redibroker")

type Config struct {
	Client *redis.Client
	EventChannels, WorkersPerChannel, EventChannelBufferSize int
}

func (cfg Config) validate() Config {
	if cfg.EventChannels <= 0 {
		cfg.EventChannels = 12
	}

	if cfg.WorkersPerChannel <= 0 {
		cfg.WorkersPerChannel = 2
	}

	if cfg.EventChannelBufferSize <= 0 {
		cfg.EventChannelBufferSize = 128
	}

	return cfg
}

type Broker struct {
	cfg			  Config
	client 		  redis.UniversalClient
	pubsub		  *redis.PubSub
	topics 		  *topicsManager
	eventChannels []chan notify.BrokerEvent
	ecIndex		  int32
	topicEvents	  chan topicEvent
	sema   		  chan struct{}
	rbp    		  *brokers.RBP
	events 		  *events.Source
	id 			  string
}

func (b *Broker) formatTopic(appID, topic string) string {
	if appID == "" && topic == "" {
		return "global"
	}
	return appID+"_"+topic
}

func (b *Broker) runWorker(ctx context.Context, evs chan notify.BrokerEvent) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-b.pubsub.Channel():
			if brokerEvent, err := b.rbp.Decode([]byte(e.Payload)); err == nil {
				if brokerEvent.BrokerID == b.id {
					continue
				}
				b.events.Emit(events.Event{Data: brokerEvent})
			} else {
				logger.Error("broker worker error:", err)
			}
		case e := <-b.topicEvents:
			if len(e.Subs) > 0 {
				err := b.pubsub.Subscribe(ctx, e.Subs...)
				if err != nil {
					logger.Error("subscribe error:", err)
				}
			}
			if len(e.Unsubs) > 0 {
				err := b.pubsub.Unsubscribe(ctx, e.Unsubs...)
				if err != nil {
					logger.Error("unsubscribe error:", err)
				}
			}
		case e := <-evs:
			if data, err := b.rbp.Encode(e); err == nil {
				switch e.Event {
				case notify.PublishEvent:
					if pubEvent, ok := e.Data.(pubsub.PublishOptions); ok {
						for _, topic := range pubEvent.Topics {
							cmd := b.client.Publish(ctx, b.formatTopic(e.AppID, topic), data)
							if cmd.Err() != nil {
								logger.Error("publish error:", cmd.Err())
							}
						}
						if len(pubEvent.Clients) > 0 || len(pubEvent.Users) > 0 {
							pubEvent.Topics = nil
							e.Data = pubEvent
							if data, err = b.rbp.Encode(e); err == nil {
								cmd := b.client.Publish(ctx, b.formatTopic(e.AppID, ""), data)
								if cmd.Err() != nil {
									logger.Error("publish error:", cmd.Err())
								}
							} else {
								logger.Error("encoding message:", err)
							}
						}
					}
				case notify.SubscribeEvent, notify.UnsubscribeEvent, notify.ConnectEvent, notify.DisconnectEvent:
					cmd := b.client.Publish(ctx, b.formatTopic(e.AppID, ""), data)
					if cmd.Err() != nil {
						logger.Error("redis publish error:", cmd.Err())
					}
				case notify.BrokerAppUpEvent:
					subErr := b.pubsub.Subscribe(ctx, b.formatTopic(e.AppID, ""))
					if subErr != nil {
						logger.Error("subscribe error:", subErr)
					}
					cmd := b.client.Publish(ctx, b.formatTopic("", ""), data)
					if cmd.Err() != nil {
						logger.Error("publish error:", cmd.Err())
					}
				case notify.BrokerAppDownEvent:
					unsubErr := b.pubsub.Unsubscribe(ctx, b.formatTopic(e.AppID, ""))
					if unsubErr != nil {
						logger.Error("publish error:", unsubErr)
					}
					cmd := b.client.Publish(ctx, b.formatTopic("", ""), data)
					if cmd.Err() != nil {
						logger.Error("publish error:", cmd.Err())
					}
				default:
					cmd := b.client.Publish(ctx, b.formatTopic("", ""), data)
					if cmd.Err() != nil {
						logger.Error("publish error:", cmd.Err())
					}
				}
			} else {
				logger.Error("encoding error:", err)
			}
		}
	}
}

func (b *Broker) lockSema() {
	b.sema <- struct{}{}
}

func (b *Broker) unlockSema() {
	select {
	case <-b.sema:
	default:
	}
}

func (b *Broker) ID() string {
	return b.ID()
}

func (b *Broker) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := b.start(ctx, wg)
		if err != nil {
			logger.Error("broker run error:", err)
		}
	}()
	wg.Wait()
}

func (b *Broker) start(ctx context.Context, wg *sync.WaitGroup) error {
	b.lockSema()
	defer b.unlockSema()

	cfg := b.cfg

	if r := b.client.Ping(ctx); r.Err() != nil {
		wg.Done()
		logger.Error("redis ping failed:", r.Err())
		return r.Err()
	}

	if b.pubsub != nil {
		if err := b.pubsub.Close(); err != nil {
			wg.Done()
			b.pubsub = nil
			logger.Error("failed to close redis pubsub:", err)
			return err
		}
	}
	b.pubsub = b.client.Subscribe(ctx, b.formatTopic("", ""))


	defer b.pubsub.Close()

	for i := 0; i < cfg.EventChannels; i++ {
		ch := make(chan notify.BrokerEvent, cfg.EventChannelBufferSize)
		b.eventChannels[i] = ch
		for j := 0; j < cfg.WorkersPerChannel; j++ {
			go b.runWorker(ctx, ch)
		}
	}
	wg.Done()
	topicsCleanTicker := time.NewTicker(TopicsCleanInterval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-topicsCleanTicker.C:
			b.topicEvents <- topicEvent{Unsubs:  b.topics.clean()}
		}
	}
}

func (b *Broker) HandlePubsubResult(appID string, r pubsub.Result) {
	b.topics.del(r.TopicsDown...)
	topics := b.topics.add(r.TopicsUp...)
	for i, topic := range topics {
		topics[i] = b.formatTopic(appID, topic)
	}

	b.topicEvents <- topicEvent{Subs: topics}
}

func (b *Broker) Handle(h func(event notify.BrokerEvent)) io.Closer {
	wrapper := func(e events.Event) {
		if e, ok := e.Data.(notify.BrokerEvent); ok && h != nil {
			h(e)
		}
	}
	return b.events.Handle(wrapper)
}

func (b *Broker) nextEventChannel() chan notify.BrokerEvent {
	index := atomic.AddInt32(&b.ecIndex, 1)
	if index >= int32(len(b.eventChannels)) {
		atomic.CompareAndSwapInt32(&b.ecIndex, index, 0)
		index = 0
	}
	return b.eventChannels[index]
}

func (b *Broker) Emit(event notify.BrokerEvent) {
	event.BrokerID = b.id
	event.Time = time.Now().UnixNano()
	b.nextEventChannel() <- event
}

func New(cfg Config) *Broker {
	cfg = cfg.validate()
	eventChannels := make([]chan notify.BrokerEvent, cfg.EventChannels)
	return &Broker{
		cfg: 		   cfg,
		client:        cfg.Client,
		topics:        newTopicsManager(),
		topicEvents:   make(chan topicEvent, 1024),
		sema:          make(chan struct{}, 1),
		rbp:           brokers.NewRBP(),
		eventChannels: eventChannels,
		events:        events.NewSource(),
		id:            uuid.New().String(),
	}
}