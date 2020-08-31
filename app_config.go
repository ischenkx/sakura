package notify

import "github.com/RomanIschenko/notify/message"

type AppConfig struct {
	ID       string
	Messages message.Storage
	PubSub   PubSubConfig
}