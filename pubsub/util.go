package pubsub

import (
	"github.com/sirupsen/logrus"
	"hash/fnv"
)

var Logger = logrus.New()

func hash(data []byte) int {
	h := fnv.New32a()
	if _, err := h.Write(data); err != nil {
		return 1
	}
	return int(h.Sum32())
}