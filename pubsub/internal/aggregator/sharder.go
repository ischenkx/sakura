package aggregator

type shardEntry struct {
	batch  Batch
	active bool
}

type Sharder struct {
	active []int
	shards []shardEntry
}

func (s *Sharder) Add(shard int, batch Batch) {
	data := &s.shards[shard]
	data.batch.Merge(batch)
	if !data.active {
		data.active = true
		s.active = append(s.active, shard)
	}
}

func (s *Sharder) AddTopics(shard int, topics ...string) {
	data := &s.shards[shard]
	if !data.active {
		data.active = true
		s.active = append(s.active, shard)
	}
	data.batch.Topics = append(data.batch.Topics, topics...)
}

func (s *Sharder) AddClients(shard int, clients ...string) {
	data := &s.shards[shard]
	if !data.active {
		data.active = true
		data.batch.Clients = append(data.batch.Clients, clients...)
		s.active = append(s.active, shard)
	}
}

func (s *Sharder) AddUsers(shard int, users ...string) {
	data := &s.shards[shard]
	if !data.active {
		data.active = true
		data.batch.Users = append(data.batch.Users, users...)
		s.active = append(s.active, shard)
	}
}

func (s *Sharder) Flush(f func(int, Batch)) {
	for _, i := range s.active {
		shard := &s.shards[i]
		shard.active = false
		f(i, shard.batch)
		shard.batch.Reset()
	}
	s.active = s.active[:0]
}

func newSharder(shards int) *Sharder {
	return &Sharder{
		active: make([]int, 0, shards),
		shards: make([]shardEntry, shards),
	}
}