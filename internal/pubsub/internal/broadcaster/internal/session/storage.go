package session

import (
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/util"
	"io"
	"sync"
)

type bucket struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func (b *bucket) delete(id string) {
	b.mu.Lock()
	delete(b.sessions, id)
	b.mu.Unlock()
}

func (b *bucket) get(id string) (*Session, bool) {
	b.mu.RLock()
	t, ok := b.sessions[id]
	b.mu.RUnlock()
	return t, ok
}

func (b *bucket) getOrCreate(id string, w io.WriteCloser, ts int64) *Session {
	b.mu.Lock()
	s, ok := b.sessions[id]
	if !ok {
		s = newSession(id, w, ts)
		b.sessions[id] = s
	}
	b.mu.Unlock()

	return s
}

func newBucket() *bucket {
	return &bucket{
		sessions: map[string]*Session{},
	}
}

type Storage struct {
	buckets []*bucket
}

func (s *Storage) bucket(id string) *bucket {
	return s.buckets[util.Hash(id) % len(s.buckets)]
}

func (s *Storage) GetOne(id string) (*Session, bool) {
	return s.bucket(id).get(id)
}

func (s *Storage) GetOrCreate(id string, w io.WriteCloser, ts int64) *Session {
	return s.bucket(id).getOrCreate(id, w, ts)
}

// TODO: use pool
func (s *Storage) Get(ids []string) []*Session {
	sessions := make([]*Session, 0, len(ids))

	for _, id := range ids {
		if sess, ok := s.bucket(id).get(id); ok {
			sessions = append(sessions, sess)
		}
	}

	return sessions
}

func (s *Storage) Delete(ids []string) {
	for _, id := range ids {
		s.bucket(id).delete(id)
	}
}

func (s *Storage) Iter(ids []string) *Iterator {
	return &Iterator{
		s:   s,
		ids: ids,
		idx: 0,
		cur: nil,
	}
}

func NewStorage(bucketsAmount int) *Storage {
	buckets := make([]*bucket, bucketsAmount)
	for i := 0; i < bucketsAmount; i++ {
		buckets[i] = newBucket()
	}
	return &Storage{buckets: buckets}
}
