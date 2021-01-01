package broadcaster

import (
	"bytes"
	"encoding/binary"
)

const (
	headerLength = 4
	buf1Op = 1
	buf2Op = 2
)

type messageWriter struct {
	buf1, buf2 *bytes.Buffer
	maxSize int
	messageBuf1, messageBuf2 []Message
}

func (w *messageWriter) Write(message Message) (bool, int) {
	if len(message.Data) + headerLength > w.maxSize {
		binary.Write(w.buf2, binary.LittleEndian, int32(len(message.Data)))
		w.buf2.Write(message.Data)
		w.messageBuf2 = append(w.messageBuf2, message)
		return true, buf2Op
	}
	if len(message.Data) + w.buf1.Len() + 4 > w.maxSize {
		return true, buf1Op
	}
	w.messageBuf1 = append(w.messageBuf1, message)
	binary.Write(w.buf1, binary.LittleEndian, int32(len(message.Data)))
	w.buf1.Write(message.Data)
	return false, buf1Op
}

func (w *messageWriter) Bytes(op int) []byte {
	switch op {
	case buf1Op:
		return w.buf1.Bytes()
	case buf2Op:
		return w.buf2.Bytes()
	}
	return nil
}

func (w *messageWriter) Reset(op int) {
	switch op {
	case buf1Op:
		w.messageBuf1 = w.messageBuf1[:0]
		w.buf1.Reset()
	case buf2Op:
		w.messageBuf2 = w.messageBuf2[:0]
		w.buf2.Reset()
	}
}

func (w *messageWriter) Messages(op int) []Message {
	switch op {
	case buf1Op:
		return w.messageBuf1
	case buf2Op:
		return w.messageBuf2
	default:
		return nil
	}
}

func newMessageWriter(maxSize int) *messageWriter {
	return &messageWriter{
		buf1:        bytes.NewBuffer(make([]byte, 0, 2<<10)),
		buf2:        bytes.NewBuffer(make([]byte, 0, 2<<15)),
		maxSize:     maxSize,
	}
}