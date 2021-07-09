package session

import (
	"github.com/ischenkx/swirl/internal/pubsub/message"
	"github.com/ischenkx/swirl/internal/pubsub/protocol"
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
	sessions []*Session
	buffer   message.Buffer
	enqueued bool
	mu 		 sync.RWMutex
}

func (g *Group) Push(messages []message.Message, guaranteedEnqueueing bool) (toBeEnqueued bool) {
	g.mu.Lock()
	g.buffer.Push(messages...)
	toBeEnqueued = !g.enqueued && g.buffer.Len() > 0
	g.enqueued = toBeEnqueued && guaranteedEnqueueing
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
	buf := message.CopyBuffer(g.buffer)
	g.buffer.Reset()
	g.mu.Unlock()
	proto.Encode(g, buf, nil)
	buf.Close()
}

func (g *Group) Add(sess *Session, ts int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	info, infoExists := g.mapper[sess.ID()]
	if infoExists && info.timestamp > ts {
		return
	}
	info.timestamp = ts
	if !infoExists || (infoExists && !info.active()) {
		g.sessions = append(g.sessions, sess)
		info.index = len(g.sessions) - 1
	}
	g.mapper[sess.ID()] = info
}

func (g *Group) Delete(id string, forced bool, ts int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	info, infoExists := g.mapper[id]

	if !infoExists {
		return
	}

	if info.timestamp > ts {
		return
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

func NewGroup() *Group {
	return &Group{
		mapper:   map[string]sessionInfo{},
		sessions: nil,
		buffer:   message.NewBuffer(),
		enqueued: false,
	}
}
