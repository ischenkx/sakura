package pubsub

import (
	"context"
	"github.com/ischenkx/swirl/internal/pubsub/changelog"
	"github.com/ischenkx/swirl/internal/pubsub/common"
	"github.com/ischenkx/swirl/internal/pubsub/engine"
	"github.com/ischenkx/swirl/internal/pubsub/message"
	"github.com/ischenkx/swirl/internal/pubsub/protocol"
	"github.com/ischenkx/swirl/internal/pubsub/subscription"
	"github.com/ischenkx/swirl/internal/pubsub/user"
	"github.com/ischenkx/swirl/internal/pubsub/util"
	"github.com/ischenkx/swirl/pkg/default/batchproto"
	"runtime"
	"time"

	"io"
)

type Config struct {
	ClientInvalidationTime time.Duration
	ProtocolProvider       protocol.Provider
	History                message.History
	Engines                int
}

func (c *Config) validate() {
	c.Engines = 100
	c.ClientInvalidationTime = time.Minute * 2
	c.ProtocolProvider = batchproto.NewProvider(1024)
}

type ConnectResult struct {
	IsReconnect         bool
	UnrecoverableTopics []string
}

type PubSub struct {
	engines          []*engine.Engine
	users            *user.Registry
	queue            chan common.Flusher
	protocolProvider protocol.Provider
	history          message.History
	config           Config
}

func (p *PubSub) engine(id string) *engine.Engine {
	idx := util.Hash(id) % len(p.engines)
	return p.engines[idx]
}

func (p *PubSub) IsActive(id string) bool {
	return p.engine(id).IsActive(id)
}

func (p *PubSub) Subscribe(id, topic string, ts int64) (changelog.ChangeLog, error) {
	cl := changelog.ChangeLog{TimeStamp: ts}
	err := p.engine(id).AddClientSubscription(id, topic, ts)
	if s, ok := p.engine(id).Session(id); ok {
		if p.engine(topic).AddTopicClient(topic, id, s, ts) {
			cl.TopicsUp = append(cl.TopicsUp, topic)
		}
	}
	return cl, err
}

func (p *PubSub) Unsubscribe(id, topic string, ts int64) (changelog.ChangeLog, error) {
	cl := changelog.ChangeLog{TimeStamp: ts}
	err := p.engine(id).DeleteClientSubscription(id, topic, ts)
	if ok, amount := p.engine(topic).DeleteTopicClient(topic, id, false, ts); ok {
		if amount == 0 {
			cl.TopicsDown = append(cl.TopicsDown, topic)
		}
	}
	return cl, err
}

func (p *PubSub) Connect(id string, w io.WriteCloser, ts int64) (ConnectResult, error) {
	var connectResult ConnectResult
	_, reconnected, disconnectTime, err := p.engine(id).UpdateClient(id, w, ts)
	if err == nil {
		p.users.UpdateClientPresence(id, w != nil, ts)
		if reconnected && p.history != nil {
			now := time.Now().UnixNano()
			for _, sub := range p.engine(id).Subscriptions(id) {
				if sub.TimeStamp > disconnectTime {
					connectResult.UnrecoverableTopics = append(connectResult.UnrecoverableTopics, sub.Topic)
					continue
				}
				if sub.Active {
					if messages, err := p.history.Snapshot(sub.Topic, disconnectTime, now); err != nil {
						connectResult.UnrecoverableTopics = append(connectResult.UnrecoverableTopics, sub.Topic)
					} else {
						p.engine(id).SendToClient(id, messages)
					}
				}
			}
		}
	}
	connectResult.IsReconnect = reconnected
	return connectResult, err
}

func (p *PubSub) Inactivate(id string, ts int64) {
	p.users.UpdateClientPresence(id, false, ts)
	p.engine(id).Inactivate(id, ts)
}

func (p *PubSub) Disconnect(id string, ts int64) changelog.ChangeLog {
	cl := changelog.ChangeLog{TimeStamp: ts}
	if info, ok := p.engine(id).Disconnect(id); ok {
		cl.ClientsDown = append(cl.ClientsDown, id)
		p.users.Delete(p.users.UserByClient(id), id, ts)
		info.Subscriptions.Iter(func(s subscription.Subscription) {
			if ok, amount := p.engine(s.Topic).DeleteTopicClient(s.Topic, id, true, ts); ok {
				if amount == 0 {
					cl.TopicsDown = append(cl.TopicsDown, s.Topic)
				}
			}
		})
	}

	return cl
}

func (p *PubSub) SendToTopic(id string, mes message.Message) {
	p.engine(id).SendToTopic(id, mes)
}

func (p *PubSub) SendToClient(id string, mes message.Message) error {
	return p.engine(id).SendToClient(id, []message.Message{mes})
}

func (p *PubSub) Subscriptions(id string) []subscription.Subscription {
	return p.engine(id).Subscriptions(id)
}

func (p *PubSub) TopicSubscribers(id string) []string {
	return p.engine(id).TopicSubscribers(id)
}

func (p *PubSub) CountTopicSubscribers(id string) int {
	return p.engine(id).CountTopicSubscribers(id)
}

func (p *PubSub) CountSubscriptions(id string) int {
	return p.engine(id).CountSubscriptions(id)
}

func (p *PubSub) Clean() []string {
	clients := make([]string, 0)
	for _, e := range p.engines {
		engineClients := e.Clean()
		timeStamp := time.Now().UnixNano()
		for _, info := range engineClients {
			clients = append(clients, info.ID)
			p.users.Delete(p.users.UserByClient(info.ID), info.ID, timeStamp)
			info.Subscriptions.Iter(func(s subscription.Subscription) {
				p.engine(s.Topic).DeleteTopicClient(s.Topic, info.ID, true, timeStamp)
			})
		}
	}
	return clients
}

func (p *PubSub) processQueue(ctx context.Context) {
	proto := p.protocolProvider.New()
	defer p.protocolProvider.Put(proto)
	for {
		select {
		case <-ctx.Done():
			return
		case f := <-p.queue:
			f.Flush(proto)
		}
	}
}

func (p *PubSub) Start(ctx context.Context) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go p.processQueue(ctx)
	}
}

func (p *PubSub) Users() *user.Registry {
	return p.users
}

func New(cfg Config) *PubSub {
	cfg.validate()
	ps := &PubSub{
		engines:          nil,
		history:          cfg.History,
		users:            user.NewRegistry(),
		queue:            make(chan common.Flusher, 8192),
		protocolProvider: cfg.ProtocolProvider,
	}

	for i := 0; i < cfg.Engines; i++ {
		ps.engines = append(ps.engines, engine.New(ps.queue, cfg.ClientInvalidationTime))
	}

	return ps
}
