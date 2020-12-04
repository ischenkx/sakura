package redibroker

import (
	"sync"
"time"
)

type channelState int

const (
	activeTopic channelState = iota
	inactiveTopic
)

type channelInfo struct {
	lastTouch int64
	state     channelState
}

type channelManager struct {
	channels map[string]*channelInfo
	inactiveChannels map[string]int64
	mu sync.RWMutex
	topicTTL time.Duration
}

func (m *channelManager) addChannels(chs []string, t int64) []string {
	added := make([]string, 0, len(chs)/4)

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ch := range chs {
		info, channelExists := m.channels[ch]
		if channelExists {
			if info.lastTouch > t {
				continue
			}
		} else {
			info =  &channelInfo{
				lastTouch: t,
				state:     activeTopic,
			}
			m.channels[ch] = info
		}

		if info.state == inactiveTopic {
			delete(m.inactiveChannels, ch)
		}
		added = append(added, ch)
	}

	return added

}

func (m *channelManager) delChannels(chs []string, t int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ch := range chs {
		info, channelExists := m.channels[ch]
		if !channelExists {
			continue
		}

		if info.lastTouch > t {
			continue
		}

		info.lastTouch = t
		info.state = inactiveTopic
		m.inactiveChannels[ch] = t
	}
}

func (m *channelManager) clean() (deletedChannels []string) {
	now := time.Now().UnixNano()
	m.mu.Lock()
	defer m.mu.Unlock()

	for ch, inactivationTime := range m.inactiveChannels {
		if now - inactivationTime > int64(time.Minute * 2) {
			delete(m.inactiveChannels, ch)
			_, channelExists := m.channels[ch]
			if !channelExists {
				continue
			}
			delete(m.channels, ch)
			deletedChannels = append(deletedChannels, ch)
		}
	}
	return
}

func newChannelManager() *channelManager {
	return &channelManager{
		channels: map[string]*channelInfo{},
		inactiveChannels: map[string]int64{},
		mu:       sync.RWMutex{},
		topicTTL: time.Minute * 2,
	}
}
