package session

import (
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/batch"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/common"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/history"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/broadcaster/internal/message"
	"io"
	"sync"
)

type Session struct {
	w io.WriteCloser
	lastUpdate int64
	enqueued bool
	id string
	histories history.PointerStorage
	buffer message.Buffer
	mu sync.RWMutex
}

func (s *Session) Update(queue chan <- common.Flusher, w io.WriteCloser, ts int64) {
	s.mu.Lock()
	if s.lastUpdate > ts {
		s.mu.Unlock()
		return
	}
	s.lastUpdate = ts
	s.w = w
	if s.w != nil {
		s.histories.Load(func(h *history.History, info history.SnapshotInfo) bool {
			if snapshot, ok := h.Load(info); ok {
				s.buffer.Push(snapshot.Slice()...)
				snapshot.Close()
			}
			return true
		})
	}
	toBeEnqueued := s.buffer.Len() > 0 && s.w != nil
	s.enqueued = toBeEnqueued
	s.mu.Unlock()
	if toBeEnqueued {
		queue <- s
	}
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Flush(batcher *batch.Batcher) {
	s.mu.Lock()
	buffer := s.buffer
	s.enqueued = false
	s.buffer = message.NewBuffer()
	s.mu.Unlock()
	defer buffer.Close()
	batcher.PutMessages(buffer)
	for {
		b, ok := batcher.Next()
		if !ok {
			break
		}
		s.Write(b)
	}
}

func (s *Session) Write(batch batch.Batch) {
	s.mu.Lock()
	if s.w == nil {
		if ptr, ok := batch.HistoryPointer(); ok {
			s.histories.Push(ptr)
		} else {
			messages := batch.Messages()
			s.buffer.Push(messages.Slice()...)
		}
	} else {
		if _, err := s.w.Write(batch.Bytes()); err != nil {
			// idk how to handle such stuff
		}
	}
	s.mu.Unlock()
}

func (s *Session) Push(q chan <- common.Flusher, messages ...message.Message) {
	s.mu.Lock()
	s.buffer.Push(messages...)
	toBeEnqueued := !s.enqueued && len(messages) > 0 && s.w != nil
	s.enqueued = toBeEnqueued
	s.mu.Unlock()

	if toBeEnqueued {
		q <- s
	}
}

func newSession(id string, w io.WriteCloser, ts int64) *Session {
	return &Session{
		w:          w,
		lastUpdate: ts,
		enqueued:   false,
		id:         id,
		histories:  history.NewPointerStorage(),
		buffer:     message.NewBuffer(),
	}
}