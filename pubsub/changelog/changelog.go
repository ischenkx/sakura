package changelog

import "time"

type Log struct {
	TopicsUp, TopicsDown,
	ClientsUp, ClientsDown,
	UsersUp, UsersDown []string
	Time int64
}

func (r *Log) Merge(r1 Log) {
	r.TopicsDown = append(r.TopicsDown, r1.TopicsDown...)
	r.TopicsUp = append(r.TopicsUp, r1.TopicsUp...)
	r.ClientsDown = append(r.ClientsDown, r1.ClientsDown...)
	r.ClientsUp = append(r.ClientsUp, r1.ClientsUp...)
	r.UsersDown = append(r.UsersDown, r1.UsersDown...)
	r.UsersUp = append(r.UsersUp, r1.UsersUp...)
}

func (r Log) Empty() bool {
	if len(r.TopicsDown) == 0 &&
		len(r.TopicsUp) == 0 &&
		len(r.ClientsUp) == 0 &&
		len(r.ClientsDown) == 0 &&
		len(r.TopicsUp) == 0 &&
		len(r.TopicsDown) == 0 {
		return true
	}
	return false
}

func New() Log {
	return Log{
		Time: time.Now().UnixNano(),
	}
}
