package memory

import (
	"github.com/ischenkx/swirl/internal/pubsub/message"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return make([]message.Message, 128)
	},
}

type topicInfo struct {
	maxlen int
	epoch  int64
	curlen int
	buf    []message.Message
	mu     sync.RWMutex
}

func (h *topicInfo) push(messages ...message.Message) {
	remlen := h.maxlen - h.curlen
	if remlen >= len(messages) {
		copy(h.buf[h.curlen:], messages)
		h.curlen += len(messages)
		return
	}
	copy(h.buf[h.curlen:h.maxlen], messages[:remlen])
	h.curlen = 0
	h.epoch += 1
	remMessages := messages[remlen:]
	h.push(remMessages...)
}

func (h *topicInfo) Push(messages ...message.Message) Pointer {
	h.mu.Lock()
	p := Pointer{
		h: h,
		info: SnapshotInfo{
			offset: h.curlen,
			epoch:  h.epoch,
		},
	}
	for _, mes := range messages {
		if opts, ok := mes.Meta.(common2.SendOptions); ok {
			if opts.DontSave {
				continue
			}
		}
		// TODO: may be inefficient (pushing a single a message)
		h.push(mes)
	}
	h.mu.Unlock()
	return p
}

func (h *topicInfo) Load(info SnapshotInfo) (message.Buffer, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	diff := h.epoch - info.epoch
	var buff message.Buffer
	if diff == 0 {
		buff = message.NewBuffer()
		messages := h.buf[info.offset:h.curlen]
		buff.Push(messages...)

		return buff, true
	} else if diff == 1 && info.offset >= h.curlen {
		buff = message.NewBuffer()
		firstPart := h.buf[info.offset:h.maxlen]
		secondPart := h.buf[:h.curlen]
		buff.Push(firstPart...)
		buff.Push(secondPart...)

		return buff, true
	}
	return buff, false
}

func (h *topicInfo) Close() {
	if h.buf == nil {
		return
	}
	bufferPool.Put(h.buf)
	h.buf = nil
}

func New() *topicInfo {
	return &topicInfo{
		maxlen: 100,
		epoch:  0,
		curlen: 0,
		buf:    bufferPool.Get().([]message.Message),
	}
}