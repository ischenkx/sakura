package broadcaster

import (
	"context"
	"sakura"
	"sakura/event"
	"sakura/util"
	"sync"
)

type Broadcaster struct {
	sakura   *sakura.Sakura
	listener event.Listener
	mu       sync.RWMutex
}

func (broadcaster *Broadcaster) Run(ctx context.Context) {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()
	broadcaster.listener = broadcaster.sakura.Events().Listener()

	util.BlockOn(ctx)
}

func (broadcaster *Broadcaster) processEvent(ev event.TopicEvent) error {
	return nil
}

func (broadcaster *Broadcaster) getEvents() {}
