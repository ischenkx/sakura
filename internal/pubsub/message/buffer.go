package message

import (
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return make([]Message, 10)
	},
}

type Buffer struct {
	messages []Message
	pooled   bool
}

func (b *Buffer) Len() int {
	return len(b.messages)
}

func (b *Buffer) Push(messages ...Message) {
	b.messages = append(b.messages, messages...)
}

func (b *Buffer) Reset() {
	b.messages = b.messages[:0]
}

func (b *Buffer) Close() {
	if b.messages == nil {
		return
	}
	if b.pooled {
		bufferPool.Put(b.messages[:0])
	}
	b.messages = nil
}

func (b Buffer) Slice() []Message {
	return b.messages
}

func BufferFrom(messages []Message) Buffer {
	return Buffer{
		messages: messages,
		pooled:   false,
	}
}

func CopyBuffer(b Buffer) Buffer {
	messages := bufferPool.Get().([]Message)[:0]

	messages = append(messages, b.messages...)

	return Buffer{
		messages: messages,
		pooled:   true,
	}
}

func NewBuffer() Buffer {
	return Buffer{
		messages: bufferPool.Get().([]Message)[:0],
		pooled:   true,
	}
}
