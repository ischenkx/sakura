package message

type History interface {
	Push(topic string, mes Message)
	Snapshot(topic string, from, to int64) ([]Message, error)
}
