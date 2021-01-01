package pubsub

import "github.com/RomanIschenko/notify/pubsub/internal/broadcaster"

type Metrics struct {
	BroadcasterMetrics broadcaster.Metrics
	Clients, Users, Topics, InactiveClients int
}
