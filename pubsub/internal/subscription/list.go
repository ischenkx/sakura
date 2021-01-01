package subscription

import "errors"

var SubscriptionNotFoundErr = errors.New("list does not contain a subscription requested")

// not supposed for concurrent use
// (requires lock)
type List struct {
	data map[string]Subscription
}

func (l List) Add(topic string, seq int64) error {
	sub := l.data[topic]
	err := sub.Activate(seq)
	if err != nil {
		return err
	}
	l.data[topic] = sub
	return nil
}

func (l List) Remove(topic string, seq int64) error {
	sub, ok := l.data[topic]
	if !ok {
		return SubscriptionNotFoundErr
	}
	err := sub.Deactivate(seq)
	if err != nil {
		return err
	}
	l.data[topic] = sub
	return nil
}

func (l List) Map() map[string]Subscription {
	return l.data
}

func NewList() List {
	return List{data: map[string]Subscription{}}
}