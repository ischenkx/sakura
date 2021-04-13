package batch

import (
	"bytes"
	"encoding/binary"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/history"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/message"
	"io"
)

const (
	bigBufferWriteOp = iota + 1001
	smallBufferWriteOp
)

type Batch struct {
	bts      []byte
	messages message.Buffer
	hptr     history.Pointer
	hptrSet  bool
}

func (b Batch) Bytes() []byte {
	return b.bts
}

func (b Batch) Messages() message.Buffer {
	return b.messages
}

func (b Batch) HistoryPointer() (history.Pointer, bool) {
	return b.hptr, b.hptrSet
}

type Batcher struct {
	messagesToWrite, shortMessages []message.Message
	smallBuffer, bigBuffer         *bytes.Buffer
	hptr                           history.Pointer
	hptrSet                        bool
	previousOp                     int
	offset                         int
	maxSize                        int
}

func (w *Batcher) encodeMessage(to io.Writer, mes message.Message) {
	binary.Write(to, binary.LittleEndian, int32(len(mes.Data)))
	to.Write(mes.Data)
}

func (w *Batcher) PutMessages(messages message.Buffer) {
	w.messagesToWrite = append(w.messagesToWrite, messages.Slice()...)
}

func (w *Batcher) PutHistoryPointer(pointer history.Pointer) {
	w.hptr = pointer
}

func (w *Batcher) Reset() {
	w.messagesToWrite = w.messagesToWrite[:0]
	w.shortMessages = nil
	w.hptrSet = false
	w.offset = 0
	w.hptr = history.Pointer{}
	w.smallBuffer.Reset()
	w.bigBuffer.Reset()
	w.previousOp = -1
}

func (w *Batcher) resetPreviousOp() {
	switch w.previousOp {
	case smallBufferWriteOp:
		w.smallBuffer.Reset()
		w.shortMessages = w.shortMessages[:0]
		w.previousOp = -1
	case bigBufferWriteOp:
		w.bigBuffer.Reset()
		w.previousOp = -1
	}
}

func (w *Batcher) Next() (b Batch, dataAvailable bool) {
	w.resetPreviousOp()
	b.hptr = w.hptr
	b.hptrSet = w.hptrSet
	mes2write := w.messagesToWrite[w.offset:]

	if len(mes2write) > 0 {
		mes := mes2write[0]
		if len(mes.Data)+4 > w.maxSize {

			w.previousOp = bigBufferWriteOp
			w.offset++
			w.encodeMessage(w.bigBuffer, mes)
			b.messages = message.BufferFrom([]message.Message{mes})
			b.bts = w.bigBuffer.Bytes()
			dataAvailable = true
			return
		}
		if len(mes.Data)+4+w.smallBuffer.Len() > w.maxSize {
			w.previousOp = smallBufferWriteOp
			b.messages = message.BufferFrom(w.shortMessages)
			b.bts = w.smallBuffer.Bytes()
			dataAvailable = true
			return
		}
		w.previousOp = -1
		w.shortMessages = append(w.shortMessages, mes)
		w.encodeMessage(w.smallBuffer, mes)
		w.offset++
		return w.Next()
	}
	if len(w.shortMessages) > 0 {
		w.previousOp = smallBufferWriteOp
		b.messages = message.BufferFrom(w.shortMessages)
		b.bts = w.smallBuffer.Bytes()
		dataAvailable = true
		return
	}
	return
}

func NewBatcher(maxSize int) *Batcher {
	return &Batcher{
		smallBuffer: bytes.NewBuffer(make([]byte, 0, 2<<10)),
		bigBuffer:   bytes.NewBuffer(make([]byte, 0, 2<<18)),
		previousOp:  -1,
		maxSize:     maxSize,
	}
}
