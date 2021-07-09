package subscription

type List struct {
	count int
	subs  map[string]Subscription
}

func (r *List) Add(topic string, ts int64) bool {
	s, ok := r.subs[topic]
	if !ok {
		s = Subscription{
			Topic:     topic,
			TimeStamp: ts,
			Active:    true,
		}
		r.subs[topic] = s
		r.count++
	}
	if s.TimeStamp >= ts {
		return s.Active
	}
	r.subs[topic] = Subscription{
		Topic:     topic,
		TimeStamp: ts,
		Active:    true,
	}
	r.count++

	return true
}

func (r *List) Delete(topic string, ts int64) bool {
	s, ok := r.subs[topic]
	if !ok {
		s = Subscription{
			Topic:     topic,
			TimeStamp: ts,
			Active:    false,
		}
		r.subs[topic] = s

		// i don't need to count--, bc subscription has never existed before
	}
	if s.TimeStamp >= ts {
		return !s.Active
	}
	r.subs[topic] = Subscription{
		Topic:     topic,
		TimeStamp: ts,
		Active:    false,
	}
	r.count--

	return true
}

func (r *List) Count() int {
	return r.count
}

func (r *List) Iter(f func(s Subscription)) {
	if f == nil {
		return
	}
	for _, subscription := range r.subs {
		f(subscription)
	}
}

func NewList() *List {
	return &List{
		count: 0,
		subs:  map[string]Subscription{},
	}
}
