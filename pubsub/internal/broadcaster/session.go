package broadcaster

import (
	"io"
	"sync"
)

type session struct {
	descriptor uint64
	writer io.WriteCloser
	buffer []Message
	failedMessages []Message
	enqueued bool
	publishQueue chan flusher
	retryEnqueued bool
	retryQueue chan uint64
	maxRetries int
	maxBufferSize int
	mu sync.Mutex
}

func (s *session) Push(messages ...Message) {
	s.mu.Lock()
	if s.writer == nil {
		for _, m := range messages {
			if len(s.buffer) >= s.maxBufferSize && s.maxBufferSize >= 0 {
				break
			}
			if m.NoBuffering {
				continue
			}
			s.buffer = append(s.buffer, m)
		}
		s.mu.Unlock()
		return
	}
	s.buffer = append(s.buffer, messages...)
	channel := s.publishQueue
	enqueued := s.enqueued
	s.enqueued = true
	s.mu.Unlock()
	if !enqueued {
		channel <- s
	}
}

func (s *session) UpdateWriter(w io.WriteCloser) {
	s.mu.Lock()
	if s.writer != nil {
		s.writer.Close()
	}
	s.writer = w
	needsPushing := len(s.buffer) > 0 && s.writer != nil && !s.enqueued
	queue := s.publishQueue
	s.enqueued = needsPushing
	s.mu.Unlock()
	if needsPushing {
		queue <- s
	}
}

func (s *session) Flush(w *messageWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.writer == nil {
		return
	}
	buffer := s.buffer
	s.buffer = nil
	s.enqueued=false
	for i := 0; i < len(buffer); i++ {
		mes := buffer[i]
		if ok, op := w.Write(mes); ok || i == len(buffer)-1 {
			messages := messageBuffer{
				bts:      w.Bytes(op),
				messages: w.Messages(op),
			}
			s.write(messages)
			if i != len(buffer)-1 {
				i--
			}
			w.Reset(op)
		}
	}
}



func (s *session) FlushFailedMessages() []Message {
	// not implemented
	return nil
}

func (s *session) write(b messageBuffer)  {
	if s.writer == nil {
		for _, m := range b.Messages() {
			if len(s.buffer) >= s.maxBufferSize && s.maxBufferSize >= 0 {
				break
			}
			if m.NoBuffering {
				continue
			}
			s.buffer = append(s.buffer, m)
		}
		return
	}
	if _, err := s.writer.Write(b.Bytes()); err != nil {
		for _, mes := range b.Messages() {
			mes.retries += 1
			if mes.retries <= s.maxRetries {
				s.failedMessages = append(s.failedMessages, mes)
			}
		}
	}
	if len(s.failedMessages) > 0 && !s.retryEnqueued {
		retryQueue := s.retryQueue
		desc := s.descriptor
		go func() {
			retryQueue <- desc
		}()
	}

}

func (s *session) Write(b messageBuffer) {
	s.mu.Lock()
	s.write(b)
	s.mu.Unlock()
}