package notify

type User struct {
	id string
	app *App
}

func (u User) ID() string {
	return u.id
}

func (u User) Clients() ([]string, error) {
	return u.app.pubsub.UserClients(u.id)
}

func (u User) Subscriptions() ([]string, error) {
	return u.app.pubsub.UserSubscriptions(u.id)
}

func (u User) Subscribe(topic string, opts ...interface{}) error {
	subOpts := SubscribeUserOptions{
		ID:    u.ID(),
		Topic: topic,
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case MetaInfoOption:
			subOpts.Meta = o.Data
		case TimeStampOption:
			subOpts.TimeStamp = o.UnixTime
		}
	}

	return u.app.SubscribeUser(subOpts)
}

func (u User) Unsubscribe(topic string, opts ...interface{}) error {
	unsubOpts := UnsubscribeUserOptions{
		ID:    u.ID(),
		Topic: topic,
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case MetaInfoOption:
			unsubOpts.Meta = o.Data
		case TimeStampOption:
			unsubOpts.TimeStamp = o.UnixTime
		}
	}

	return u.app.UnsubscribeUser(unsubOpts)
}

func (u User) UnsubscribeAll(opts ...interface{}) error {
	unsubOpts := UnsubscribeUserOptions{
		ID:  u.ID(),
		All: true,
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case MetaInfoOption:
			unsubOpts.Meta = o.Data
		case TimeStampOption:
			unsubOpts.TimeStamp = o.UnixTime
		}
	}
	return u.app.UnsubscribeUser(unsubOpts)
}

func (u User) Emit(name string, data ...interface{}) {
	ev := newEvent(name, data)
	ev.Users = []string{u.id}
	u.app.Emit(ev)
}
