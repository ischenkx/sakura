package broadcaster

import (
	"context"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
)

type Config struct {
	MaxBatchSize                            int
	ClientBufferSize                        int
	SessionsPerSubGroup                     int
	PublishQueueCap, RetryQueueCap, Buckets int
}

func (c *Config) validate() {
	if c.MaxBatchSize <= 0 {
		c.MaxBatchSize = 2<<15
	}
	if c.Buckets <= 0 {
		c.Buckets = 2<<9
	}
	if c.PublishQueueCap <= 0 {
		c.PublishQueueCap = 2<<16
	}
	if c.RetryQueueCap <= 0 {
		c.RetryQueueCap = 2<<16
	}
	if c.SessionsPerSubGroup <= 0 {
		c.SessionsPerSubGroup = 2<<7
	}
}

type bucket struct {
	sessions map[uint64]*session
	groups map[uint64]*group
	pushers map[uint64]pusher
	publishQueue chan flusher
	retryQueue chan uint64
	config Config
	mu sync.RWMutex
}

func (b *bucket) join(desc uint64, s *session) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if g, ok := b.groups[desc]; ok {
		g.Add(s)
	}
}

func (b *bucket) leave(desc uint64, sd uint64) (length int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if g, ok := b.groups[desc]; ok {
		length = g.Delete(sd)
		if length == 0 {
			delete(b.groups, desc)
		}
	}
	return
}

func (b *bucket) newGroup(desc uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.groups[desc]; ok {
		return
	}
	g := &group{
		descriptor:          desc,
		publishQueue: b.publishQueue,
		subs: map[uint64]int{},
		sessionsPerSubGroup: b.config.SessionsPerSubGroup,
	}
	b.groups[desc] = g
	b.pushers[desc] = g
}

func (b *bucket) deleteGroup(desc uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.groups[desc]; ok {
		delete(b.pushers, desc)
		delete(b.groups, desc)
	}
}

func (b *bucket) deleteSession(desc uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.sessions[desc]; ok {
		delete(b.pushers, desc)
		delete(b.sessions, desc)
	}
}

func (b *bucket) updateSession(desc uint64, w io.WriteCloser) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if s, ok := b.sessions[desc]; ok {
		go s.UpdateWriter(w)
		return
	}
	s := &session{
		descriptor:     desc,
		writer:         w,
		retryQueue:     b.retryQueue,
		publishQueue: b.publishQueue,
		maxBufferSize: b.config.ClientBufferSize,
	}
	b.sessions[desc] = s
	b.pushers[desc] = s
}

func (b *bucket) checkGroupValidity(desc uint64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if g, ok := b.groups[desc]; ok {
		if g.Len() == 0 {
			delete(b.groups, desc)
			return false
		}
	}
	return true
}

func (b *bucket) push(desc uint64, messages []Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if w, ok := b.pushers[desc]; ok {
		go w.Push(messages...)
	}
}

func (b *bucket) getSession(desc uint64) (*session, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	s, ok := b.sessions[desc]
	return s, ok
}

type Broadcaster struct {
	buckets []*bucket
	publishQueue chan flusher
	config Config
	desc uint64
	retryQueue chan uint64
}

func (b *Broadcaster) nextDescriptor() uint64 {
	return atomic.AddUint64(&b.desc, 1)
}

func (b *Broadcaster) bucket(desc uint64) *bucket {
	return b.buckets[int(desc % uint64(len(b.buckets)))]
}

func (b *Broadcaster) UpdateSession(desc uint64, w io.WriteCloser) {
	b.bucket(desc).updateSession(desc, w)
}

// if group is empty it is deleted
func (b *Broadcaster) CheckGroupValidity(desc uint64) bool {
	return b.bucket(desc).checkGroupValidity(desc)
}

func (b *Broadcaster) NewSession(w io.WriteCloser) uint64 {
	desc := b.nextDescriptor()
	b.bucket(desc).updateSession(desc, w)
	return desc
}

func (b *Broadcaster) DeleteSession(desc uint64) {
	b.bucket(desc).deleteSession(desc)
}

func (b *Broadcaster) NewGroup() uint64 {
	desc := b.nextDescriptor()
	b.bucket(desc).newGroup(desc)
	return desc
}

func (b *Broadcaster) DeleteGroup(desc uint64) {
	b.bucket(desc).deleteGroup(desc)
}

func (b *Broadcaster) Leave(sd, gd uint64) int {
	return b.bucket(gd).leave(gd, sd)
}

func (b *Broadcaster) Join(sd, gd uint64) {
	if s, ok := b.bucket(sd).getSession(sd); ok {
		b.bucket(gd).join(gd, s)
	}
}

func (b *Broadcaster) Push(desc uint64, messages []Message) {
	b.bucket(desc).push(desc, messages)
}

func (b *Broadcaster) startWorker(ctx context.Context) {
	mw := newMessageWriter(b.config.MaxBatchSize)
	for {
		select {
		case <-ctx.Done():
			return
		case f := <-b.publishQueue:
			f.Flush(mw)
		case <-b.retryQueue:
			// dunno what to do yet
		}
	}
}



func (b *Broadcaster) Start(ctx context.Context) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go b.startWorker(ctx)
	}
}

func (b *Broadcaster) Metrics() (m Metrics) {
	for _, buck := range b.buckets {
		buck.mu.RLock()
		m.Sessions += len(buck.sessions)
		m.Groups += len(buck.groups)
		buck.mu.RUnlock()
	}
	m.PushersEnqueued = len(b.publishQueue)
	m.RetriesEnqueued = len(b.retryQueue)
	return m
}

func New(cfg Config) *Broadcaster {
	cfg.validate()
	b := &Broadcaster{
		buckets:      make([]*bucket, cfg.Buckets),
		publishQueue: make(chan flusher, cfg.PublishQueueCap),
		retryQueue: make(chan uint64, cfg.RetryQueueCap),
		config: cfg,
	}

	for i := 0; i < len(b.buckets); i++ {
		b.buckets[i] = &bucket{
			sessions: map[uint64]*session{},
			groups: map[uint64]*group{},
			pushers: map[uint64]pusher{},
			publishQueue: b.publishQueue,
			retryQueue:   b.retryQueue,
			config: cfg,
		}
	}
	return b
}