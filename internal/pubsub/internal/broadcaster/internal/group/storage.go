package group

import (
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/session"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/util"
	"sync"
)

type bucket struct {
	groups map[string]*Group
	mu     sync.RWMutex
}

func (b *bucket) get(id string) (*Group, bool) {
	b.mu.RLock()
	t, groupExists := b.groups[id]
	b.mu.RUnlock()
	return t, groupExists
}

func (b *bucket) getOrCreate(id string) *Group {
	t, groupExists := b.get(id)

	if !groupExists {
		b.mu.Lock()
		t, groupExists = b.groups[id]
		if !groupExists {
			t = newGroup()
			b.groups[id] = t
		}
		b.mu.Unlock()
	}

	return t
}

func (b *bucket) join(id string, ts int64, sessions []*session.Session) {
	if len(sessions) == 0 {
		return
	}
	t := b.getOrCreate(id)
	t.add(ts, sessions)
}

func (b *bucket) leave(id string, ts int64, sessions []string, forced bool) {
	if len(sessions) == 0 {
		return
	}
	if t, ok := b.get(id); ok {
		t.delete(ts, sessions, forced)
	}
}

func newBucket() *bucket {
	return &bucket{
		groups: map[string]*Group{},
	}
}

type Storage struct {
	buckets []*bucket
}

func (s *Storage) bucket(id string) *bucket {
	return s.buckets[util.Hash(id)%len(s.buckets)]
}

func (s *Storage) Join(id string, ts int64, clients []*session.Session) {
	s.bucket(id).join(id, ts, clients)
}

func (s *Storage) Leave(id string, ts int64, clients []string, forced bool) {
	s.bucket(id).leave(id, ts, clients, forced)
}

func (s *Storage) Get(id string) (*Group, bool) {

	return s.bucket(id).get(id)
}

func NewStorage(bucketsAmount int) *Storage {
	buckets := make([]*bucket, bucketsAmount)
	for i := 0; i < bucketsAmount; i++ {
		buckets[i] = newBucket()
	}
	return &Storage{buckets: buckets}
}
