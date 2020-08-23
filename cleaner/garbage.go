package cleaner

import (
	"sync"
	"time"
)

type Garbage struct {
	data map[string]time.Time
	expirationTime time.Duration
	mu sync.Mutex
}

func (g *Garbage) Add(d string) {
	g.mu.Lock()
	g.data[d] = time.Now()
	g.mu.Unlock()
}

func (g *Garbage) Del(d string) {
	g.mu.Lock()
	delete(g.data, d)
	g.mu.Unlock()
}

func (g *Garbage) Flush() []string {
	now := time.Now()
	flushedData := make([]string, 0)
	g.mu.Lock()
	expTime := g.expirationTime
	for d, t := range g.data {
		if now.Sub(t) >= expTime {
			flushedData = append(flushedData, d)
			delete(g.data, d)
		}
	}
	g.mu.Unlock()
	return flushedData
}

func NewGarbage(expTime time.Duration) *Garbage {
	return &Garbage{
		data: map[string]time.Time{},
		expirationTime: expTime,
	}
}