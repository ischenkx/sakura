package notify

type AppConfig struct {
	ID string
	Messages MessageStorage
	PubSub PubSubConfig
}