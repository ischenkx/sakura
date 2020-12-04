package distributor

import (
	"github.com/RomanIschenko/notify/internal"
	"github.com/RomanIschenko/notify/pubsub/changelog"
	"github.com/RomanIschenko/notify/pubsub/client_id"
)

type Distributor struct {
	registry topicsRegistry
	sharderPool *sharderPool
	shards int
}

func (d *Distributor) TopicsAmount() int {
	return d.registry.Amount()
}

func (d *Distributor) AggregateSharded(b Batch, f func(int, Batch) changelog.Log) (r changelog.Log) {
	sharder := d.sharderPool.Get()
	for _, clientID := range b.Clients {
		id := clientid.ID(clientID)
		hash, err := id.Hash()
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
	d.registry.loadTopics(b.Topics, sharder)
	sharder.Flush(func(s int, b Batch) {
		tmpRes := f(s, b)
		d.registry.processResult(s, tmpRes)
		r.Merge(tmpRes)
	})
	d.sharderPool.Put(sharder)
	return
}

func (d *Distributor) ProcessResult(s int, r changelog.Log) {
	d.registry.processResult(s, r)
}

func (d *Distributor) Aggregate(f func(int) changelog.Log) (r changelog.Log) {
	for i := 0; i < d.shards; i++ {
		res := f(i)
		r.Merge(res)
		d.registry.processResult(i, res)
	}
	return
}

func (d *Distributor) Topics() []string {
	topics := make([]string, 0, d.TopicsAmount())

	for _, bucket := range d.registry {
		bucket.mu.RLock()
		for id := range bucket.b {
			topics = append(topics, id)
		}
		bucket.mu.RUnlock()
	}

	return topics
}

func New(shards int) *Distributor {
	return &Distributor{
		registry:    newTopicsRegistry(16),
		shards:      shards,
		sharderPool: newSharderPool(shards),
	}
}