package broadcaster

import (
	"context"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/batch"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/common"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/group"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/message"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/session"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/mutator"
	"log"
	"runtime"
)

type Broadcaster struct {
	topics, users *group.Storage
	sessions      *session.Storage
	queue         chan common.Flusher
}

func (b *Broadcaster) commit(m *Mutator) {
	m.Iter(func(data interface{}) {
		switch mut := data.(type) {
		case mutator.UserDetachedMutation:
			b.users.Leave(mut.User, m.Timestamp(), mut.Clients, mut.Forced)
		case mutator.UserAttachedMutation:
			b.users.Join(mut.User, m.Timestamp(), b.sessions.Get(mut.Clients))
		case mutator.ClientUpdatedMutation:
			b.sessions.
				GetOrCreate(mut.Client, mut.Writer, m.Timestamp()).
				Update(b.queue, mut.Writer, m.Timestamp())
		case mutator.ClientsDeletedMutation:
			b.sessions.Delete(mut.Clients)
		case mutator.TopicSubscribedMutation:
			b.topics.Join(mut.Topic, m.Timestamp(), b.sessions.Get(mut.Clients))
		case mutator.TopicUnsubscribedMutation:
			b.topics.Leave(mut.Topic, m.Timestamp(), mut.Clients, mut.Forced)
		}
	})
}

func (b *Broadcaster) Mutator(ts int64) Mutator {
	return Mutator{
		Mutator:   mutatorsPool.Get().(mutator.Mutator),
		closed:    false,
		timestamp: ts,
		b:         b,
	}
}

func (b *Broadcaster) Broadcast(sessions, users, topics []string, data []byte) {

	mes := message.New(data)

	for _, id := range sessions {
		if s, ok := b.sessions.GetOne(id); ok {
			s.Push(b.queue, mes)
		} else {
			log.Println("failed to get session:", id)
		}
	}

	for _, id := range users {
		if u, ok := b.users.Get(id); ok {
			u.Push(b.queue, mes)
		}
	}

	for _, id := range topics {
		if t, ok := b.topics.Get(id); ok {
			t.Push(b.queue, mes)
		}
	}
}

func (b *Broadcaster) processFlushers(ctx context.Context) {
	batcher := batch.NewBatcher(8192)
	for {
		select {
		case <-ctx.Done():
			return
		case f := <-b.queue:
			f.Flush(batcher)
		}
		batcher.Reset()
	}
}

func (b *Broadcaster) Start(ctx context.Context) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go b.processFlushers(ctx)
	}
}

func New() *Broadcaster {
	return &Broadcaster{
		topics:   group.NewStorage(1000),
		users:    group.NewStorage(1000),
		sessions: session.NewStorage(1000),
		queue:    make(chan common.Flusher, 100000),
	}
}
