package topic

type Topic struct {
	clients map[string]struct{}
}

func (t *Topic) Add(id string) {
	t.clients[id] = struct{}{}
}

func (t *Topic) Del(id string) {
	delete(t.clients, id)
}

func (t *Topic) Len() int {
	return len(t.clients)
}

func (t *Topic) Clients() map[string]struct{} {
	return t.clients
}

func New() *Topic {
	return &Topic{clients: map[string]struct{}{}}
}