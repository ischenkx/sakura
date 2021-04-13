package notify

type ActionBuilder struct {
	app *App
	clients, users, topics []string
	metaInfo interface{}
	timeStamp int64
}

func (a ActionBuilder) WithTimeStamp(t int64) ActionBuilder {
	a.timeStamp = t
	return a
}

func (a ActionBuilder) WithMetaInfo(data interface{}) ActionBuilder {
	a.metaInfo = data
	return a
}

func (a ActionBuilder) WithClients(clients ...string) ActionBuilder {
	a.clients = clients
	return a
}

func (a ActionBuilder) WithUsers(users ...string) ActionBuilder {
	a.users = users
	return a
}

func (a ActionBuilder) WithTopics(topics ...string) ActionBuilder {
	a.topics = topics
	return a
}

func (a ActionBuilder) Subscribe() ActionBuilder {
	a.app.Subscribe(SubscribeOptions{
		Clients:   a.clients,
		Users:     a.users,
		Topics:    a.topics,
		TimeStamp: a.timeStamp,
		Meta:      a.metaInfo,
	})
	return a
}

func (a ActionBuilder) Unsubscribe() ActionBuilder {
	a.app.Unsubscribe(UnsubscribeOptions{
		Clients:   a.clients,
		Users:     a.users,
		Topics:    a.topics,
		TimeStamp: a.timeStamp,
		Meta:      a.metaInfo,
	})
	return a
}

func (a ActionBuilder) UnsubscribeAll() ActionBuilder {
	a.app.Unsubscribe(UnsubscribeOptions{
		All:          true,
		Topics:       a.topics,
		TimeStamp:    a.timeStamp,
		Meta:         a.metaInfo,
	})
	return a
}

func (a ActionBuilder) UnsubscribeAllFromTopic() ActionBuilder {
	a.app.Unsubscribe(UnsubscribeOptions{
		AllFromTopic: true,
		Topics:       a.topics,
		TimeStamp:    a.timeStamp,
		Meta:         a.metaInfo,
	})
	return a
}

func (a ActionBuilder) Publish(data []byte) ActionBuilder {
	a.app.Publish(PublishOptions{
		Clients:   a.clients,
		Users:     a.users,
		Topics:    a.topics,
		Data:      data,
		TimeStamp: a.timeStamp,
		Meta:      a.metaInfo,
	})
	return a
}