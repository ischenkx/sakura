package pubsub

type TopicChangeLog interface {
	TopicID() string
	ClientsAdded() []string
	UsersAdded() []string
	ClientsDeleted() []string
	UsersDeleted() []string
}

type ChangeLog interface {
	TopicsCreated() []string
	TopicsDeleted()      []string
	UsersCreated()       []string
	UsersDeleted()       []string
	ClientsCreated()     []string
	ClientsDeleted()     []string
	ClientsInactivated() []string
	Topics()             []TopicChangeLog
	Timestamp()          int64
}
