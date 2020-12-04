package internal

import "hash/fnv"

func Hash(data []byte) int {
	h := fnv.New32a()
	if _, err := h.Write(data); err != nil {
		return 1
	}
	return int(h.Sum32())
}
