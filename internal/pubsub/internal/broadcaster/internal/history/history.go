package history

import (
	"fmt"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/message"

	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return make([]message.Message, 0, 2048)
	},
}

type History struct {
	maxlen int
	epoch int64
	curlen int
	buf []message.Message
	mu sync.RWMutex
}

func (h *History) push(messages ...message.Message) {
	remlen := h.maxlen - h.curlen
	if remlen >= len(messages) {
		h.curlen += len(messages)
		copy(h.buf[h.curlen:], messages)
		return
	}
	copy(h.buf, messages[:remlen])
	h.curlen = 0
	h.epoch += 1
	remMessages := messages[remlen:]
	h.push(remMessages...)
}

func (h *History) Push(messages ...message.Message) Pointer {

	h.mu.Lock()
	h.push(messages...)
	if h.curlen > h.maxlen {
		fmt.Println("EROROROROROROR:", len(h.buf))
	}
	p := Pointer{
		h:    h,
		info: SnapshotInfo{
			offset: h.curlen,
			epoch:  h.epoch,
		},
	}
	h.mu.Unlock()

	return p
}

func (h *History) Load(info SnapshotInfo) (message.Buffer, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	diff := h.epoch - info.epoch

	var buff message.Buffer

	if diff == 0 {
		buff = message.NewBuffer()
		messages := h.buf[info.offset:h.curlen]
		buff.Push(messages...)
		return buff, true
	} else if diff == 1 && info.offset > h.curlen {
		buff = message.NewBuffer()
		firstPart := h.buf[info.offset:]
		secondPart := h.buf[:h.curlen]
		buff.Push(firstPart...)
		buff.Push(secondPart...)
		return buff, true
	}
	return buff, false
}

func (h *History) Close() {
	if h.buf == nil {
		return
	}
	bufferPool.Put(h.buf)
	h.buf = nil
}

var count = int64(0)

func New() *History {
	return &History{
		maxlen: 3,
		epoch:  0,
		curlen: 0,
		buf:    make([]message.Message, 4096),
	}
}