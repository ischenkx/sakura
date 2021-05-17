package history

import (
	common2 "github.com/ischenkx/notify/internal/pubsub/common"
	message2 "github.com/ischenkx/notify/internal/pubsub/message"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return make([]message2.Message, 128)
	},
}


type History struct {
	maxlen int
	epoch  int64
	curlen int
	buf    []message2.Message
	mu     sync.RWMutex
}

func (h *History) push(messages ...message2.Message) {
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

func (h *History) Push(messages ...message2.Message) Pointer {
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

func (h *History) Load(info SnapshotInfo) (message2.Buffer, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	diff := h.epoch - info.epoch
	var buff message2.Buffer
	if diff == 0 {
		buff = message2.NewBuffer()
		messages := h.buf[info.offset:h.curlen]
		buff.Push(messages...)

		return buff, true
	} else if diff == 1 && info.offset >= h.curlen {
		buff = message2.NewBuffer()
		firstPart := h.buf[info.offset:h.maxlen]
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


func New() *History {
	return &History{
		maxlen: 100,
		epoch:  0,
		curlen: 0,
		buf:    bufferPool.Get().([]message2.Message),
	}
}
