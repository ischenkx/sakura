package pubsub

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"github.com/RomanIschenko/notify/pubsub/clientid"
	"github.com/RomanIschenko/notify/pubsub/internal/aggregator"
	"github.com/RomanIschenko/notify/pubsub/namespace"
	"time"
)

const DefaultShardsAmount = 128
const DefaultCleanInterval = time.Minute*3

type Config struct {
	Shards 		  int
	ClientConfig   ClientConfig
	CleanInterval time.Duration
}

func (cfg *Config) validate() {
	if cfg.Shards <= 0 {
		cfg.Shards = DefaultShardsAmount
	}

	if cfg.CleanInterval <= 0 {
		cfg.CleanInterval = DefaultCleanInterval
	}
}

type Pubsub struct {
	shards     []*shard
	config     Config
	events     *eventsRegistry
	nsRegistry *namespace.Registry
	aggregator *aggregator.Registry
	proxy      *proxyRegistry
}

func (p *Pubsub) Proxy(ctx context.Context) *Proxy {
	return p.proxy.hub(ctx)
}

func (p *Pubsub) Events(ctx context.Context) *EventsHub {
	return p.events.hub(ctx)
}

func (p *Pubsub) Metrics() Metrics {
	metrics := &Metrics{
		Topics: p.aggregator.TopicsAmount(),
	}
	for _, shard := range p.shards {
		sm := shard.Metrics()
		metrics.Merge(&sm)
	}
	return *metrics
}

func (p *Pubsub) NamespaceRegistry() *namespace.Registry {
	return p.nsRegistry
}

func (p *Pubsub) Clients() []string {
	var clients []string

	for _, s := range p.shards {
		clients = append(clients, s.Clients()...)
	}

	return clients
}

func (p *Pubsub) Users() []string {
	var users []string

	for _, s := range p.shards {
		users = append(users, s.Users()...)
	}

	return users
}

func (p *Pubsub) Topics() []string {
	return p.aggregator.Topics()
}

func  (p *Pubsub) Publish(opts PublishOptions) error {
	opts.validate()
	b := aggregator.Batch{Clients: opts.Clients, Users: opts.Users, Topics: opts.Topics}
	p.proxy.emitPublish(&opts)
	p.aggregator.AggregateSharded(b, func(shardIdx int, b aggregator.Batch) changelog.Log {
		shard := p.shards[shardIdx]
		opts.Topics = b.Topics
		opts.Users = b.Users
		opts.Clients = b.Clients
		shard.Publish(opts)
		return changelog.Log{}
	})
	p.events.emitPublish(opts)
	return nil
}

func (p *Pubsub) Subscribe(opts SubscribeOptions) changelog.Log {
	opts.validate()
	b := aggregator.Batch{Clients: opts.Clients, Users: opts.Users}
	p.proxy.emitSubscribe(&opts)
	res := p.aggregator.AggregateSharded(b, func(shardIdx int, b aggregator.Batch) changelog.Log {
		shard := p.shards[shardIdx]
		opts.Users = b.Users
		opts.Clients = b.Clients
		return shard.Subscribe(opts)
	})
	p.events.emitSubscribe(opts, res)
	return res
}

func (p *Pubsub) Unsubscribe(opts UnsubscribeOptions) changelog.Log {
	opts.validate()
	b := aggregator.Batch{Clients: opts.Clients, Users: opts.Users}
	p.proxy.emitUnsubscribe(&opts)
	res := p.aggregator.AggregateSharded(b, func(shardIdx int, b aggregator.Batch) changelog.Log {
		shard := p.shards[shardIdx]
		opts.Topics = b.Topics
		opts.Users = b.Users
		opts.Clients = b.Clients
		return shard.Unsubscribe(opts)
	})
	p.events.emitUnsubscribe(opts, res)
	return res
}

func (p *Pubsub) Connect(opts ConnectOptions) (*Client, error) {
	if opts.ID == "" {
		opts.ID = clientid.New("")
	}
	h, err := clientid.Hash(opts.ID)
	if err != nil {
		return nil, err
	}
	shard := p.shards[h % len(p.shards)]
	p.proxy.emitConnect(&opts)
	c, res, err := shard.Connect(opts)
	p.events.emitConnect(opts, c, res)
	return c, err
}

func (p *Pubsub) Disconnect(opts DisconnectOptions) changelog.Log {
	opts.validate()
	b := aggregator.Batch{Clients: opts.Clients, Users: opts.Users}
	p.proxy.emitDisconnect(&opts)
	res := p.aggregator.AggregateSharded(b, func(shardIdx int, b aggregator.Batch) changelog.Log {
		shard := p.shards[shardIdx]
		opts.Users = b.Users
		opts.Clients = b.Clients
		return shard.Disconnect(opts)
	})
	p.events.emitDisconnect(opts, res)
	return res
}

func (p *Pubsub) InactivateClient(client *Client) {
	if client == nil {
		return
	}
	h := client.Hash()
	idx := h % len(p.shards)
	shard := p.shards[idx]
	p.proxy.emitInactivate(client)
	res := shard.InactivateClient(client)
	p.aggregator.ProcessResult(idx, res)
	p.events.emitInactivate(client.ID())
}

func (p *Pubsub) Clean() {
	res := p.aggregator.Aggregate(func(s int) changelog.Log {
		shard := p.shards[s]
		return shard.Clean()
	})
	if !res.Empty() {
		p.events.emitChange(res)
	}
}

// StartCleaner cleans the whole pubsub system
// each Config.CleanInterval
// It is a blocking call
func (p *Pubsub) StartCleaner(ctx context.Context) {
	cleaner := time.NewTicker(p.config.CleanInterval)
	defer cleaner.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-cleaner.C:
			p.Clean()
		}
	}
}

func New(config Config) *Pubsub {
	config.validate()
	shards := make([]*shard, config.Shards)
	nsRegistry := namespace.NewRegistry()

	for i := range shards {
		shards[i] = newShard(nsRegistry, config.ClientConfig)
	}

	return &Pubsub{
		shards:     shards,
		config:     config,
		proxy:      newProxyRegistry(),
		events:     newEventsRegistry(),
		nsRegistry: nsRegistry,
		aggregator: aggregator.New(config.Shards),
	}
}