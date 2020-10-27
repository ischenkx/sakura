package pubsub

import (
	"hash/fnv"
)

func hash(data []byte) int {
	h := fnv.New32a()
	if _, err := h.Write(data); err != nil {
		return 1
	}
	return int(h.Sum32())
}