package subscription

import "errors"

type List struct {
	data map[string]Subscription
}

func (l List) Add(id string, ts int64) error {
	s := l.data[id]
	if s.TimeStamp > ts {
		return errors.New("invalid timestamp")
	}
	s.TimeStamp = ts
	s.Active = true
	l.data[id] = s
	return nil
}

func (l List) Delete(id string, ts int64) error {
	s, ok := l.data[id]
	if !ok {
		return errors.New("no subscription")
	}
	s.TimeStamp = ts
	s.Active = false
	l.data[id] = s
	return nil
}

func (l List) Map() map[string]Subscription {
	return l.data
}

func NewList() List {
	return List{map[string]Subscription{}}
}
