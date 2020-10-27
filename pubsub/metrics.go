package pubsub

type Metrics struct {
	Users, Clients, Topics int
}

func (m *Metrics) Merge(otherMetrics *Metrics) {
	m.Clients += otherMetrics.Clients
	m.Users += otherMetrics.Users
	m.Topics += otherMetrics.Topics
}

