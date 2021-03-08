package pubsub

type topic struct {
	clients map[string]struct{}
}

func (t *topic) Add(id string) {
	t.clients[id] = struct{}{}
}

func (t *topic) Del(id string) {
	delete(t.clients, id)
}

func (t *topic) Len() int {
	return len(t.clients)
}

func newTopic() *topic {
	return &topic{clients: map[string]struct{}{}}
}