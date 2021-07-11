package memory

import "github.com/ischenkx/swirl/internal/pubsub/message"

type History struct {

}

func (h *History) Push(topic string, mes message.Message) {

}

func (h *History) Snapshot(topic string, from, to int64) ([]message.Message, error) {

}

