package registry

import (
	"github.com/RomanIschenko/notify/internal"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"github.com/RomanIschenko/notify/pubsub/clientid"
)

type Registry struct {
	topics      topicsRegistry
	sharderPool *sharderPool
	shards      int
}

func (d *Registry) TopicsAmount() int {
	return d.topics.Amount()
}

func (d *Registry) AggregateSharded(b Batch, f func(int, Batch) changelog.Log) (r changelog.Log) {
	sharder := d.sharderPool.Get()
	for _, clientID := range b.Clients {
		hash, err := clientid.Hash(clientID)
		if err != nil {
			continue
		}
		shard := hash % d.shards
		sharder.AddClients(shard, clientID)
	}
	for _, userID := range b.Users {
		hash := internal.Hash([]byte(userID))
		shard := hash % d.shards
		sharder.AddClients(shard, userID)
	}
	d.topics.loadTopics(b.Topics, sharder)
	sharder.Flush(func(s int, b Batch) {
		tmpRes := f(s, b)
		d.topics.processResult(s, tmpRes)
		r.Merge(tmpRes)
	})
	d.sharderPool.Put(sharder)
	return
}

func (d *Registry) ProcessResult(s int, r changelog.Log) {
	d.topics.processResult(s, r)
}

func (d *Registry) Aggregate(f func(int) changelog.Log) (r changelog.Log) {
	for i := 0; i < d.shards; i++ {
		res := f(i)
		r.Merge(res)
		d.topics.processResult(i, res)
	}
	return
}

func (d *Registry) Topics() []string {
	topics := make([]string, 0, d.TopicsAmount())

	for _, bucket := range d.topics {
		bucket.mu.RLock()
		for id := range bucket.b {
			topics = append(topics, id)
		}
		bucket.mu.RUnlock()
	}

	return topics
}

func New(shards int) *Registry {
	return &Registry{
		topics:      newTopicsRegistry(16),
		shards:      shards,
		sharderPool: newSharderPool(shards),
	}
}