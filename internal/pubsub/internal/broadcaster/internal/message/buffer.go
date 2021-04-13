package message

import (
	"sync"
	"sync/atomic"
)

var created = int64(0)
var poolput = int64(0)
var closed = int64(0)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		//fmt.Println("buffers:", atomic.AddInt64(&created, 1), atomic.LoadInt64(&closed), atomic.LoadInt64(&poolput))
		return make([]Message, 0, 256)
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
	atomic.AddInt64(&closed, 1)
	if b.messages == nil {
		return
	}
	if b.pooled {
		atomic.AddInt64(&poolput, 1)
		bufferPool.Put(b.messages[:0])
	}
	b.messages = nil
}

func (b *Buffer) Slice() []Message {
	return b.messages
}

func BufferFrom(messages []Message) Buffer {
	return Buffer{
		messages: messages,
		pooled:   false,
	}
}

func NewBuffer() Buffer {
	return Buffer{
		messages: bufferPool.Get().([]Message),
		pooled:   true,
	}
}
