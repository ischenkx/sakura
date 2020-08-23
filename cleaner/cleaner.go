package cleaner

import (
	"sync"
	"time"
)

type Cleanable interface {
	Clean()
}

type CleanerState int

const (
	RunningCleaner CleanerState = iota + 1
	StoppedCleaner
)

type Cleaner struct {
	cleanable Cleanable
	interval  time.Duration
	state     CleanerState
	closer    chan struct{}
	mu        sync.RWMutex
}

func (cleaner *Cleaner) SetCleanable(c Cleanable) {
	if c == nil {
		return
	}
	cleaner.mu.Lock()
	defer cleaner.mu.Unlock()
	cleaner.cleanable = c
}

func (cleaner *Cleaner) SetInterval(i time.Duration) {
	cleaner.mu.RLock()
	defer cleaner.mu.RUnlock()
	cleaner.interval = i
}

func (cleaner *Cleaner) startCleaning(closer chan struct{}) {
	for {
		cleaner.mu.RLock()
		cleanable := cleaner.cleanable
		interval := cleaner.interval
		cleaner.mu.RUnlock()
		select {
			case <-closer:
				return
			case <-time.After(interval):
				if cleanable != nil {
					cleanable.Clean()
				}
		}
	}
}

func (cleaner *Cleaner) Run() {
	cleaner.mu.Lock()
	if cleaner.state == StoppedCleaner {
		closer := make(chan struct{})
		cleaner.closer = closer
		cleaner.state = RunningCleaner
		go cleaner.startCleaning(closer)
	}
	cleaner.mu.Unlock()
}

func (cleaner *Cleaner) Stop() {
	cleaner.mu.Lock()
	if cleaner.state == RunningCleaner {
		cleaner.closer <- struct{}{}
		cleaner.state = StoppedCleaner
	}
	cleaner.mu.Unlock()
}

func New(cleanable Cleanable, interval time.Duration) *Cleaner {
	return &Cleaner{
		cleanable: cleanable,
		interval:  interval,
		state:     StoppedCleaner,
	}
}