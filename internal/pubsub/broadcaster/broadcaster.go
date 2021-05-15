package broadcaster

import (
	"context"
	common2 "github.com/RomanIschenko/notify/internal/pubsub/broadcaster/internal/common"
	group2 "github.com/RomanIschenko/notify/internal/pubsub/broadcaster/internal/group"
	session2 "github.com/RomanIschenko/notify/internal/pubsub/broadcaster/internal/session"
	message2 "github.com/RomanIschenko/notify/internal/pubsub/message"
	mutator2 "github.com/RomanIschenko/notify/internal/pubsub/mutator"
	protocol2 "github.com/RomanIschenko/notify/internal/pubsub/protocol"
	"log"
	"runtime"
)

type Broadcaster struct {
	topics, users *group2.Storage
	sessions      *session2.Storage

	protocolProvider protocol2.Provider

	queue         chan common2.Flusher
}

func (b *Broadcaster) commit(m *Mutator) {
	m.Iter(func(data interface{}) {
		switch mut := data.(type) {
		case mutator2.UserDetachedMutation:
			b.users.Leave(mut.User, m.Timestamp(), mut.Clients, mut.Forced)
		case mutator2.UserAttachedMutation:
			b.users.Join(mut.User, m.Timestamp(), b.sessions.Get(mut.Clients))
		case mutator2.ClientUpdatedMutation:
			ses := b.sessions.GetOrCreate(mut.Client, mut.Writer, m.Timestamp())
			toBeEnqueued := ses.Update(mut.Writer, m.Timestamp(), true)
			if toBeEnqueued {
				b.queue <- ses
			}
		case mutator2.ClientsDeletedMutation:
			b.sessions.Delete(mut.Clients)
		case mutator2.TopicSubscribedMutation:
			b.topics.Join(mut.Topic, m.Timestamp(), b.sessions.Get(mut.Clients))
		case mutator2.TopicUnsubscribedMutation:
			b.topics.Leave(mut.Topic, m.Timestamp(), mut.Clients, mut.Forced)
		}
	})
}

func (b *Broadcaster) Mutator(ts int64) Mutator {
	return Mutator{
		Mutator:   mutatorsPool.Get().(mutator2.Mutator),
		closed:    false,
		timestamp: ts,
		b:         b,
	}
}

func (b *Broadcaster) Broadcast(sessions, users, topics []string, messages []message2.Message) {

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
			log.Println("failed to get topic!!!")
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

func New(protocolProvider protocol2.Provider) *Broadcaster {
	return &Broadcaster{
		topics:           group2.NewStorage(1000),
		users:            group2.NewStorage(1000),
		sessions:         session2.NewStorage(1000),
		queue:            make(chan common2.Flusher, 100000),
		protocolProvider: protocolProvider,
	}
}
