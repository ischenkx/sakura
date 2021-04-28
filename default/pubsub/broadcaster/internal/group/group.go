package group

import (
	"github.com/RomanIschenko/notify/default/pubsub/broadcaster/internal/history"
	"github.com/RomanIschenko/notify/default/pubsub/broadcaster/internal/session"
	"github.com/RomanIschenko/notify/pubsub/message"
	"github.com/RomanIschenko/notify/pubsub/protocol"
	"log"
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
	sessions []*session.Session
	buffer   message.Buffer
	history  *history.History
	enqueued bool

	mu sync.RWMutex
}

func (g *Group) Push(messages []message.Message, guaranteedEnqueueing bool) (toBeEnqueued bool) {
	g.mu.Lock()
	g.buffer.Push(messages...)
	toBeEnqueued = !g.enqueued && g.buffer.Len() > 0
	g.enqueued = toBeEnqueued&&guaranteedEnqueueing
	g.mu.Unlock()
	return
}

func (g *Group) Write(batch message.Batch, meta interface{}) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, sess := range g.sessions {
		sess.Write(batch, meta)
	}
}

func (g *Group) Flush(proto protocol.Protocol) {
	g.mu.Lock()
	g.enqueued = false
	hptr := g.history.Push(g.buffer.Slice()...)
	buf := message.CopyBuffer(g.buffer)
	g.buffer.Reset()
	g.mu.Unlock()
	proto.Encode(g, buf, hptr)
	buf.Close()
}

func (g *Group) add(ts int64, cls []*session.Session) {
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
			log.Println("forcibly deleting:", id)
			delete(g.mapper, id)
		}
	}
}

func newGroup() *Group {
	return &Group{
		mapper:   map[string]sessionInfo{},
		sessions: make([]*session.Session, 0, 200),
		buffer:   message.NewBuffer(),
		history:  history.New(),
		enqueued: false,
	}
}
