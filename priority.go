package notify

type Priority int64

const (
	PluginPriority Priority = iota + 1
	UserPriority
)
