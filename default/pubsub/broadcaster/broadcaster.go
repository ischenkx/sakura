package broadcaster

import (
	"context"
	"github.com/RomanIschenko/notify/default/pubsub/broadcaster/internal/common"
	"github.com/RomanIschenko/notify/default/pubsub/broadcaster/internal/group"
	"github.com/RomanIschenko/notify/default/pubsub/broadcaster/internal/session"
	"github.com/RomanIschenko/notify/default/pubsub/mutator"
	"github.com/RomanIschenko/notify/pubsub/message"
	"github.com/RomanIschenko/notify/pubsub/protocol"
	"runtime"
)

type Broadcaster struct {
	topics, users *group.Storage
	sessions      *session.Storage

	protocolProvider protocol.Provider

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
			ses := b.sessions.GetOrCreate(mut.Client, mut.Writer, m.Timestamp())
			toBeEnqueued := ses.Update(mut.Writer, m.Timestamp(), true)
			if toBeEnqueued {
				b.queue <- ses
			}
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

func (b *Broadcaster) Broadcast(sessions, users, topics []string, messages []message.Message) {

	for _, id := range sessions {
		c, ok := b.sessions.GetOne(id)
		if !ok {
			continue
		}
		toBeEnqueued := c.Push(messages, true)
		if toBeEnqueued {
			b.queue <- c
		}
	}

	for _, id := range users {
		u, ok := b.users.Get(id)
		if !ok {
			continue
		}
		toBeEnqueued := u.Push(messages, true)
		if toBeEnqueued {
			b.queue <- u
		}
	}

	for _, id := range topics {
		t, ok := b.topics.Get(id)
		if !ok {
			continue
		}
		toBeEnqueued := t.Push(messages, true)
		if toBeEnqueued {
			b.queue <- t
		}
	}
}

func (b *Broadcaster) processQueue(ctx context.Context) {
	proto := b.protocolProvider.New()
	defer b.protocolProvider.Put(proto)
	for {
		select {
		case <-ctx.Done():
			return
		case f := <-b.queue:
			f.Flush(proto)
		}
	}
}

func (b *Broadcaster) Start(ctx context.Context) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go b.processQueue(ctx)
	}
}

func New(protocolProvider protocol.Provider) *Broadcaster {
	return &Broadcaster{
		topics:           group.NewStorage(1000),
		users:            group.NewStorage(1000),
		sessions:         session.NewStorage(1000),
		queue:            make(chan common.Flusher, 100000),
		protocolProvider: protocolProvider,
	}
}
