package pubsub

import "time"

type ChangeLog struct {
	ClientsUp, ClientsDown,
	UsersUp, UsersDown,
	TopicsUp, TopicsDown []string
	Time int64
}

func (c *ChangeLog) setTime() {
	c.Time = time.Now().UnixNano()
}