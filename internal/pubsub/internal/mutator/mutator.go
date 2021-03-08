package mutator

import (
	"io"
	"sync"
)

var stringsPool = &sync.Pool{
	New: func() interface{} {
		return make([]string, 0, 128)
	},
}

type MutationType int

type Mutation struct {
	Type MutationType
	Value interface{}
}

type Mutator interface {
	Subscribe(session string, topic string)
	Unsubscribe(session string, topic string, forced bool)
	AttachUser(session string, user string)
	DetachUser(session string, user string, forced bool)
	UpdateClient(session string, w io.WriteCloser)
	DeleteClient(session string)
	Iter(func(interface{}))
	Reset()
}

type mutator struct {
	subscribedTopics           map[string][]string
	unsubscribedTopics         map[string][]string
	attachedUsers              map[string][]string
	detachedUsers              map[string][]string
	forciblyDetachedUsers      map[string][]string
	forciblyUnsubscribedTopics map[string][]string
	updatedClients             map[string]io.WriteCloser
	deletedClients             []string
}

func (m *mutator) pushToKey(hm map[string][]string, key, data string) {
	if arr, ok := hm[key]; ok {
		hm[key] = append(arr, data)
	} else {
		arr = stringsPool.Get().([]string)
		arr = append(arr, data)
		hm[key] = arr
	}
}

func (m *mutator) Subscribe(session string, topic string) {
	m.pushToKey(m.subscribedTopics, topic, session)
}

func (m *mutator) Unsubscribe(session string, topic string, forced bool) {
	hm := m.unsubscribedTopics
	if forced {
		hm = m.forciblyUnsubscribedTopics
	}
	m.pushToKey(hm, topic, session)
}

func (m *mutator) AttachUser(session string, user string) {
	m.pushToKey(m.attachedUsers, user, session)
}

func (m *mutator) DetachUser(session string, user string, forced bool) {
	hm := m.detachedUsers
	if forced {
		hm = m.forciblyDetachedUsers
	}
	m.pushToKey(hm, user, session)
}

func (m *mutator) UpdateClient(session string, w io.WriteCloser) {
	m.updatedClients[session] = w
}

func (m *mutator) DeleteClient(session string) {
	m.deletedClients = append(m.deletedClients, session)
}

func (m *mutator) Iter(f func(interface{})) {
	if f == nil {
		return
	}

	if len(m.deletedClients) > 0 {
		f(ClientsDeletedMutation{Clients: m.deletedClients})
	}

	for s, w := range m.updatedClients {
		f(ClientUpdatedMutation{
			Client: s,
			Writer: w,
		})
	}

	for s, ids := range m.attachedUsers {
		f(UserAttachedMutation{
			User: s,
			Clients: ids,
		})
	}

	for s, ids := range m.detachedUsers {
		f(UserDetachedMutation{
			User: s,
			Clients: ids,
		})
	}

	for s, ids := range m.forciblyDetachedUsers {
		f(UserDetachedMutation{
			User: s,
			Clients: ids,
			Forced: true,
		})
	}

	for s, ids := range m.subscribedTopics {
		f(TopicSubscribedMutation{
			Topic: s,
			Clients: ids,
		})
	}

	for s, ids := range m.unsubscribedTopics {
		f(TopicUnsubscribedMutation{
			Topic: s,
			Clients: ids,
		})
	}

	for s, ids := range m.forciblyUnsubscribedTopics {
		f(TopicUnsubscribedMutation{
			Topic: s,
			Clients: ids,
			Forced: true,
		})
	}
}

func (m *mutator) Reset() {
	m.deletedClients = m.deletedClients[:0]

	for s, _ := range m.updatedClients {
		delete(m.updatedClients, s)
	}

	for s, ids := range m.attachedUsers {

		delete(m.attachedUsers, s)
		stringsPool.Put(ids[:0])
	}

	for s, ids := range m.detachedUsers {
		delete(m.detachedUsers, s)
		stringsPool.Put(ids[:0])
	}

	for s, ids := range m.subscribedTopics {
		delete(m.subscribedTopics, s)
		stringsPool.Put(ids[:0])
	}

	for s, ids := range m.unsubscribedTopics {
		delete(m.unsubscribedTopics, s)
		stringsPool.Put(ids[:0])
	}

	for s, ids := range m.forciblyDetachedUsers {
		delete(m.forciblyDetachedUsers, s)
		stringsPool.Put(ids[:0])
	}

	for s, ids := range m.forciblyUnsubscribedTopics {
		delete(m.forciblyUnsubscribedTopics, s)
		stringsPool.Put(ids[:0])
	}
}

func New() Mutator {
	return &mutator{
		subscribedTopics:           map[string][]string{},
		unsubscribedTopics:         map[string][]string{},
		attachedUsers:              map[string][]string{},
		detachedUsers:              map[string][]string{},
		forciblyDetachedUsers: map[string][]string{},
		forciblyUnsubscribedTopics: map[string][]string{},
		updatedClients:             map[string]io.WriteCloser{},
		deletedClients:             make([]string, 128),
	}
}

