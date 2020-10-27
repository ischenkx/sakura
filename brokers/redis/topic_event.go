package redibroker

type topicEvent struct {
	AppID string
	Subs []string
	Unsubs []string
}
