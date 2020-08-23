package notify

import (
	"sync"
)

type ChannelState int

const (
	ActiveChannel ChannelState = iota + 1
	InactiveChannel
)

type Channel struct {
	members map[string]*Client
	id string
	state ChannelState
	mu sync.RWMutex
}

func (ch *Channel) tryInactivate() bool {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if len(ch.members) == 0 {
		ch.state = InactiveChannel
		return true
	}
	return false
}

func (ch *Channel) len() int {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	return len(ch.members)
}

func (ch *Channel) add(client *Client) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if ch.state != ActiveChannel {
		return
	}
	if ch.members == nil {
		ch.members = map[string]*Client{}
	}
	ch.members[client.id] = client
}

func (ch *Channel) del(id string) (int, bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	deleted := false
	_, deleted = ch.members[id]
	delete(ch.members, id)
	l := len(ch.members)

	return l, deleted
}

func (ch *Channel) send(mes Message) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	if ch.state != ActiveChannel {
		return
	}
	if len(ch.members) > 0 {
		for _, client := range ch.members {
			client.send(mes)
		}
	}
}
