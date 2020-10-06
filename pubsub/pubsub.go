package pubsub

import (
	"context"
	"github.com/sirupsen/logrus"
	"hash/fnv"
	"time"
)

const DefaultShardsAmount = 128
const DefaultTopicBuckets = 32
const DefaultCleanInterval = time.Minute*3

type Config struct {
	Shards 		  int
	ShardConfig   ShardConfig
	PubQueueConfig PubQueueConfig
	CleanInterval time.Duration
	Logger		  logrus.FieldLogger
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

	if cfg.Logger == nil {
		cfg.Logger = logrus.New()
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
	topics topicProvider
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

func (p *Pubsub) NS() *namespaceRegistry {
	return p.nsRegistry
}

func (p *Pubsub) processActionResult(shard int, r result) {
	if shard >= len(p.shards) || shard < 0 {
		//fmt.Println("wefojwefiojweiojwef")
		return
	}
	p.topics.add(r.topicsUp, shard)
	p.topics.del(r.topicsDown, shard)
}

func (p *Pubsub) clean() {
	for _, shard := range p.shards {
		shard.Clean()
	}
}

func (p *Pubsub) Publish(opts PublishOptions) {
	b := batch{opts.Clients, opts.Users, opts.Topics}

	p.config.Logger.Info("pub", opts)

	for shardIdx, batch := range p.distribute(b) {
		shard := p.shards[shardIdx]
		opts.Topics = batch.topics
		opts.Users = batch.users
		opts.Clients = batch.clients
		p.processActionResult(shardIdx, shard.Publish(opts))
	}
}

//func (p *Pubsub) logStats() {
//	if file, err := os.OpenFile("stats.txt", os.O_CREATE, 0666); err == nil {
//		for i, shard := range p.shards {
//			shard.mu.RLock()
//			sTopics := len(shard.topics)
//			sClients := len(shard.clients)
//			sUsers := len(shard.users)
//			shard.mu.RUnlock()
//			data := fmt.Sprintf(
//				"SHARD %v:\n\tTopics:%v,\n\tClients:%v,\n\tUsers:%v\n",
//				i, sTopics, sClients, sUsers,
//				)
//			file.Write([]byte(data))
//		}
//	}
//
//}

func (p *Pubsub) Subscribe(opts SubscribeOptions) {
	b := batch{opts.Clients, opts.Users, nil}

	p.config.Logger.Info("sub", opts)

	for shardIdx, batch := range p.distribute(b) {
		shard := p.shards[shardIdx]
		opts.Users = batch.users
		opts.Clients = batch.clients
		p.processActionResult(shardIdx, shard.Subscribe(opts))
	}
}

func (p *Pubsub) Unsubscribe(opts UnsubscribeOptions) {
	b := batch{opts.Clients, opts.Users, opts.Topics}

	p.config.Logger.Info("unsub", opts)

	for shardIdx, batch := range p.distribute(b) {
		shard := p.shards[shardIdx]
		opts.Topics = batch.topics
		opts.Users = batch.users
		opts.Clients = batch.clients
		p.processActionResult(shardIdx, shard.Unsubscribe(opts))
	}
}

func (p *Pubsub) Connect(opts ConnectOptions) (*Client, error) {

	p.config.Logger.Info("con", opts.ID)

	h, err := opts.ID.Hash()
	if err != nil {
		return nil, err
	}
	shard := p.shards[h % len(p.shards)]
	return shard.Connect(opts)
}

func (p *Pubsub) Disconnect(opts DisconnectOptions) {
	b := batch{opts.Clients, opts.Users, nil}

	p.config.Logger.Info("discon", opts)

	for shardIdx, batch := range p.distribute(b) {
		shard := p.shards[shardIdx]
		opts.Users = batch.users
		opts.Clients = batch.clients
		p.processActionResult(shardIdx, shard.Disconnect(opts))
	}
}

func (p *Pubsub) Start(ctx context.Context) {
	cleaner := time.NewTicker(p.config.CleanInterval)
	p.queue.Start(ctx)

	p.config.Logger.Info("pubsub started")

	for {
		select {
		case <-ctx.Done():
			p.config.Logger.Info("pubsub is done")
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

	p.config.Logger.Info("client_inactivation", client.ID())

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
	topicProvider := newTopicProvider(config.TopicBuckets)

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