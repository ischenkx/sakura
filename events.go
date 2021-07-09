package swirl

type EmitOptions struct {
	Clients, Users, Topics []string
	EventOptions           EventOptions
}

type AppEvents interface {
	OnEmit(f func(*App, EmitOptions))
	OnEvent(f func(*App, Client, EventOptions))
	OnConnect(f func(*App, ConnectOptions, Client))
	OnDisconnect(f func(*App, Client))
	OnReconnect(f func(*App, ConnectOptions, Client))
	OnInactivate(f func(*App, Client))
	OnError(f func(*App, error))
	OnChange(f func(*App, ChangeLog))
	OnClientSubscribe(f func(*App, Client, string, int64))
	OnUserSubscribe(f func(*App, User, string, int64))
	OnClientUnsubscribe(f func(*App, Client, string, int64))
	OnUserUnsubscribe(f func(*App, User, string, int64))
	Close()
}

type ClientEvents interface {
	OnEmit(func(EventOptions))
	OnEvent(func(EventOptions))
	OnUnsubscribe(func(topic string, ts int64))
	OnSubscribe(func(topic string, ts int64))
	OnReconnect(func(int64))
	OnDisconnect(func(ts int64))
	Close()
}

type UserEvents interface {
	OnEmit(func(EventOptions))
	OnEvent(func(EventOptions))
	OnUnsubscribe(func(topic string, ts int64))
	OnSubscribe(func(topic string, ts int64))
	OnClientConnect(func(Client, int64))
	OnClientReconnect(func(Client, int64))
	OnClientDisconnect(func(id string, ts int64))
	Close()
}

type TopicEvents interface {
	OnClientSubscribe(func(Client))
	OnUserSubscribe(func(User))
	OnClientUnsubscribe(func(Client))
	OnUserUnsubscribe(func(User))
	OnEmit(func(options EventOptions))
	Close()
}
