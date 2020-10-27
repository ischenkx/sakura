package pubsub

import (
	"context"
	"hash/fnv"
	"time"
)

const DefaultShardsAmount = 128
const DefaultTopicBuckets = 32
const DefaultCleanInterval = time.Minute*3

type ActionResult struct {
	TopicsUp []string
	TopicsDown []string
}

type Config struct {
	Shards 		  int
	ShardConfig   ShardConfig
	PubQueueConfig PubQueueConfig
	CleanInterval time.Duration
	TopicBuckets  int
}

func (cfg Config) validate() Config {
	if cfg.Shards <= 0 {
		cfg.Shards = DefaultShardsAmount
	}

	if cfg.TopicBuckets <= 0 {
		cfg.TopicBuckets = DefaultTopicBuckets
	}

	if cfg.CleanInterval <= 0 {
		cfg.CleanInterval = DefaultCleanInterval
	}

	if cfg.TopicBuckets <= 0 {
		cfg.TopicBuckets = 8
	}

	return cfg
}

type batch struct {
	clients, users, topics []string
}

type Pubsub struct {
	shards []*shard
	config Config
	queue pubQueue
	nsRegistry *namespaceRegistry
	topics topicDispatcher
}

func (p *Pubsub) hash(b []byte) (int, error) {
	h := fnv.New32a()
	if _, err := h.Write(b); err != nil {
		return 0, err
	}
	return int(h.Sum32()), nil
}

func (p *Pubsub) distribute(b batch) map[int]batch {
	db := map[int]batch{}
	l := len(p.shards)

	for _, clientID := range b.clients {
		id := ClientID(clientID)
		if h, err := id.Hash(); err == nil {
			idx := h % l
			b1 := db[idx]
			b1.clients = append(b1.clients, clientID)
			db[idx] = b1
		}
	}

	for _, userID := range b.users {
		h := hash([]byte(userID))
		idx := h % l
		b1 := db[idx]
		b1.users = append(b1.users, userID)
		db[idx] = b1
	}

	for s, topics := range p.topics.get(b.topics) {
		b1 := db[s]
		b1.topics = topics
		db[s] = b1
	}

	return db
}

func (p *Pubsub) clean() {
	for _, shard := range p.shards {
		shard.Clean()
	}
}

func (p *Pubsub) processShardResults(r shardResultCollector) (res Result) {
	res.TopicsUp, res.TopicsDown = p.topics.processShardResults(r)
	return
}

func (p *Pubsub) Metrics() Metrics {
	metrics := &Metrics{
		Topics: p.topics.Amount(),
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

func (p *Pubsub) Publish(opts PublishOptions) Result {
	b := batch{opts.Clients, opts.Users, opts.Topics}

	c := newShardResultCollector()

	for shardIdx, batch := range p.distribute(b) {
		shard := p.shards[shardIdx]
		opts.Topics = batch.topics
		opts.Users = batch.users
		opts.Clients = batch.clients
		c.collect(shardIdx, shard.Publish(opts))
	}

	return p.processShardResults(c)
}

func (p *Pubsub) Subscribe(opts SubscribeOptions) Result {
	b := batch{opts.Clients, opts.Users, nil}

	c := newShardResultCollector()

	for shardIdx, batch := range p.distribute(b) {
		shard := p.shards[shardIdx]
		opts.Users = batch.users
		opts.Clients = batch.clients
		c.collect(shardIdx, shard.Subscribe(opts))
	}

	return p.processShardResults(c)
}

func (p *Pubsub) Unsubscribe(opts UnsubscribeOptions) Result {
	b := batch{opts.Clients, opts.Users, opts.Topics}

	c := newShardResultCollector()

	for shardIdx, batch := range p.distribute(b) {
		shard := p.shards[shardIdx]
		opts.Topics = batch.topics
		opts.Users = batch.users
		opts.Clients = batch.clients
		c.collect(shardIdx, shard.Unsubscribe(opts))
	}

	return p.processShardResults(c)
}

func (p *Pubsub) Connect(opts ConnectOptions) (*Client, error) {

	h, err := opts.ID.Hash()
	if err != nil {
		return nil, err
	}
	shard := p.shards[h % len(p.shards)]
	return shard.Connect(opts)
}

func (p *Pubsub) Disconnect(opts DisconnectOptions) Result {
	b := batch{opts.Clients, opts.Users, nil}

	c := newShardResultCollector()

	for shardIdx, batch := range p.distribute(b) {
		shard := p.shards[shardIdx]
		opts.Users = batch.users
		opts.Clients = batch.clients
		c.collect(shardIdx, shard.Disconnect(opts))
	}

	return p.processShardResults(c)
}

func (p *Pubsub) Start(ctx context.Context) {
	p.queue.Start(ctx)

	cleaner := time.NewTicker(p.config.CleanInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-cleaner.C:
			p.clean()
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
	shard.InactivateClient(client)
}

func New(config Config) *Pubsub {
	config = config.validate()
	shards := make([]*shard, config.Shards)

	queue := newPubQueue(config.PubQueueConfig)
	nsRegistry := newNamespaceRegistry()
	topicProvider := newTopicDispatcher(config.TopicBuckets)

	for i := range shards {
		shards[i] = newShard(queue, nsRegistry, config.ShardConfig)
	}

	return &Pubsub{
		shards: shards,
		queue: queue,
		config: config,
		nsRegistry: nsRegistry,
		topics: topicProvider,
	}
}