package mutator

import "io"

type UserDetachedMutation struct {
	Clients []string
	User string
	Forced bool
}

type UserAttachedMutation struct {
	Clients []string
	User string
}

type ClientUpdatedMutation struct {
	Client string
	Writer io.WriteCloser
}

type ClientsDeletedMutation struct {
	Clients []string
}

type TopicSubscribedMutation struct {
	Clients []string
	Topic string
}

type TopicUnsubscribedMutation struct {
	Clients []string
	Topic string
	Forced bool
}

type ForcedDeletionMutation struct {
	Clients []string
}