package pubsub

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"github.com/RomanIschenko/notify/pubsub/internal/registry"
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
	events 	   *eventsHub
	nsRegistry *namespaceRegistry
	distrib *registry.Registry
}

func (p *Pubsub) Events(ctx context.Context) *ContextEvents {
	return p.events.ctxHub(ctx)
}

func (p *Pubsub) Clean() {
	res := p.distrib.Aggregate(func(s int) changelog.Log {
		shard := p.shards[s]
		return shard.Clean()
	})
	if !res.Empty() {
		p.events.emitChange(res)
	}
}

func (p *Pubsub) Metrics() Metrics {
	metrics := &Metrics{
		Topics: p.distrib.TopicsAmount(),
	}
	for _, shard := range p.shards {
		sm := shard.Metrics()
		metrics.Merge(&sm)
	}
	return *metrics
}

func (p *Pubsub) NS() *namespaceRegistry {
	return p.nsRegistry
}

func  (p *Pubsub) Publish(opts PublishOptions) error {
	opts.validate()
	b := registry.Batch{Clients: opts.Clients, Users: opts.Users, Topics: opts.Topics}
	p.events.emitBeforePublish(opts)
	p.distrib.AggregateSharded(b, func(shardIdx int, b registry.Batch) changelog.Log {
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
	b := registry.Batch{Clients: opts.Clients, Users: opts.Users}
	res := p.distrib.AggregateSharded(b, func(shardIdx int, b registry.Batch) changelog.Log {
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
	b := registry.Batch{Clients: opts.Clients, Users: opts.Users}
	res := p.distrib.AggregateSharded(b, func(shardIdx int, b registry.Batch) changelog.Log {
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
		opts.ID = NewClientID("")
	}
	h, err := HashClientID(opts.ID)
	if err != nil {
		return nil, err
	}
	shard := p.shards[h % len(p.shards)]
	c, res, err := shard.Connect(opts)
	p.events.emitConnect(opts, c, res)
	return c, err
}

func (p *Pubsub) Disconnect(opts DisconnectOptions) changelog.Log {
	opts.validate()
	b := registry.Batch{Clients: opts.Clients, Users: opts.Users}
	res := p.distrib.AggregateSharded(b, func(shardIdx int, b registry.Batch) changelog.Log {
		shard := p.shards[shardIdx]
		opts.Users = b.Users
		opts.Clients = b.Clients
		return shard.Disconnect(opts)
	})
	p.events.emitDisconnect(opts, res)
	return res
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

func (p *Pubsub) InactivateClient(client *Client) {
	if client == nil {
		return
	}

	h := client.Hash()
	idx := h % len(p.shards)
	shard := p.shards[idx]
	res := shard.InactivateClient(client)
	p.distrib.ProcessResult(idx, res)
	p.events.emitInactivate(client.ID())
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
	return p.distrib.Topics()
}

func New(config Config) *Pubsub {
	config.validate()
	shards := make([]*shard, config.Shards)
	nsRegistry := newNamespaceRegistry()

	for i := range shards {
		shards[i] = newShard(nsRegistry, config.ClientConfig)
	}

	return &Pubsub{
		shards:     shards,
		config:     config,
		events:     newEventsHub(),
		nsRegistry: nsRegistry,
		distrib:    registry.New(config.Shards),
	}
}