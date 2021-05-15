package group

import (
	history2 "github.com/RomanIschenko/notify/internal/pubsub/broadcaster/internal/history"
	session2 "github.com/RomanIschenko/notify/internal/pubsub/broadcaster/internal/session"
	message2 "github.com/RomanIschenko/notify/internal/pubsub/message"
	protocol2 "github.com/RomanIschenko/notify/internal/pubsub/protocol"
	"sync"
)

type sessionInfo struct {
	timestamp int64
	// -1 => not in Group
	index int
}

func (info sessionInfo) active() bool {
	return info.index >= 0
}

type Group struct {
	mapper   map[string]sessionInfo
	sessions []*session2.Session
	buffer   message2.Buffer
	history  *history2.History
	enqueued bool

	mu sync.RWMutex
}

func (g *Group) Push(messages []message2.Message, guaranteedEnqueueing bool) (toBeEnqueued bool) {
	g.mu.Lock()
	g.buffer.Push(messages...)
	toBeEnqueued = !g.enqueued && g.buffer.Len() > 0
	g.enqueued = toBeEnqueued&&guaranteedEnqueueing
	g.mu.Unlock()
	return
}

func (g *Group) Write(batch message2.Batch, meta interface{}) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, sess := range g.sessions {
		sess.Write(batch, meta)
	}
}

func (g *Group) Flush(proto protocol2.Protocol) {
	g.mu.Lock()
	g.enqueued = false
	hptr := g.history.Push(g.buffer.Slice()...)
	buf := message2.CopyBuffer(g.buffer)
	g.buffer.Reset()
	g.mu.Unlock()
	proto.Encode(g, buf, hptr)
	buf.Close()
}

func (g *Group) add(ts int64, cls []*session2.Session) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, c := range cls {
		info, infoExists := g.mapper[c.ID()]
		if infoExists && info.timestamp > ts {
			continue
		}
		info.timestamp = ts
		if !infoExists || (infoExists && !info.active()) {
			g.sessions = append(g.sessions, c)
			info.index = len(g.sessions) - 1
		}
		g.mapper[c.ID()] = info
	}
}

func (g *Group) delete(ts int64, ids []string, forced bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, id := range ids {
		info, infoExists := g.mapper[id]
		if !infoExists {
			continue
		}

		if info.timestamp > ts {
			continue
		}

		if info.active() {
			last := g.sessions[len(g.sessions)-1]
			g.sessions[info.index] = last
			g.sessions = g.sessions[:len(g.sessions)-1]
			lastInfo := g.mapper[last.ID()]
			lastInfo.index = info.index
			info.index = -1
			g.mapper[last.ID()] = lastInfo
		}

		info.timestamp = ts
		g.mapper[id] = info

		if forced {
			delete(g.mapper, id)
		}
	}
}

func newGroup() *Group {
	return &Group{
		mapper:   map[string]sessionInfo{},
		sessions: make([]*session2.Session, 0, 200),
		buffer:   message2.NewBuffer(),
		history:  history2.New(),
		enqueued: false,
	}
}
