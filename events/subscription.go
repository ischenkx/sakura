package events

type Subscription struct {
	pubsub *Pubsub
	channel chan Event
}

func (s Subscription) Close() {
	s.pubsub.deleteSub(s.channel)
}

func (s Subscription) Channel() <-chan Event {
	return s.channel
}