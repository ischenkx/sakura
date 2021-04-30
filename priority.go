package notify

type Priority int64

const (
	PluginPriority Priority = 0
	UserPriority            = 1 << 16
)
