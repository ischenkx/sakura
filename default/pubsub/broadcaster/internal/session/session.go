package session

import (
	"github.com/RomanIschenko/notify/default/pubsub/broadcaster/internal/history"
	"github.com/RomanIschenko/notify/default/pubsub/common"
	"github.com/RomanIschenko/notify/pubsub/message"
	"github.com/RomanIschenko/notify/pubsub/protocol"
	"io"
	"sync"
)

type Session struct {
	w          io.WriteCloser
	lastUpdate int64
	enqueued   bool
	id         string
	histories  history.PointerStorage
	buffer     message.Buffer
	mu         sync.RWMutex
}

func (s *Session) Update(w io.WriteCloser, ts int64, guaranteedEnqueueing bool) (toBeEnqueued bool) {
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
				//fmt.Println("loaded:", snapshot.Len(), snapshot.Slice())
				s.buffer.Push(snapshot.Slice()...)
				snapshot.Close()
			}
			return true
		})
	}
	toBeEnqueued = s.buffer.Len() > 0 && s.w != nil
	s.enqueued = toBeEnqueued&&guaranteedEnqueueing
	s.mu.Unlock()
	return
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Flush(proto protocol.Protocol) {
	s.mu.Lock()
	buffer := s.buffer
	s.enqueued = false
	s.buffer = message.NewBuffer()
	s.mu.Unlock()
	proto.Encode(s, buffer, nil)
	buffer.Close()
}

func (s *Session) unsafeWrite(batch message.Batch, meta interface{}) {
	if s.w == nil {
		if ptr, ok := meta.(history.Pointer); ok {
			s.histories.Push(ptr)
		} else {
			s.buffer.Push(batch.Buffer().Slice()...)
		}
	} else {
		if _, err := s.w.Write(batch.Bytes()); err != nil {
			// idk how to handle such stuff
		}
	}
}

func (s *Session) Write(batch message.Batch, meta interface{}) {
	s.mu.Lock()
	s.unsafeWrite(batch, meta)
	s.mu.Unlock()
}

func (s *Session) Push(messages []message.Message, guaranteedEnqueueing bool) (toBeEnqueued bool) {
	s.mu.Lock()
	if s.w == nil {
		for _, mes := range messages {
			if opts, ok := mes.Meta.(common.SendOptions); ok {
				if opts.DontSave {
					continue
				}
			}
			s.buffer.Push(mes)
		}
	} else {
		s.buffer.Push(messages...)
		toBeEnqueued = !s.enqueued && len(messages) > 0
		s.enqueued = toBeEnqueued&&guaranteedEnqueueing
		s.mu.Unlock()
	}
	return
}

func (s *Session) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buffer.Close()
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
