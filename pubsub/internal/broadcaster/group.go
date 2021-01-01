package broadcaster

import (
	"fmt"
	"sync"
)

type subGroup struct {
	enqueued          bool
	publishQueue      chan flusher
	retryEnqueued     bool
	buffer            []Message
	sessions          map[uint64]*session
	intermediateState bool
	deletedSessions   []uint64
	maxSize 		  int
	mu                sync.RWMutex
}

func (sg *subGroup) Push(messages ...Message) {
	sg.mu.Lock()
	sg.buffer = append(sg.buffer, messages...)
	queue := sg.publishQueue
	needsEnqueueing := !sg.enqueued
	sg.mu.Unlock()
	if needsEnqueueing {
		queue <- sg
	}
}

func (sg *subGroup) Flush(w *messageWriter) {
	sg.mu.Lock()
	sg.intermediateState = true
	sg.enqueued = false
	buffer := sg.buffer
	sg.buffer = nil
	sg.mu.Unlock()
	for i := 0; i < len(buffer); i++ {
		mes := buffer[i]
		if ok, op := w.Write(mes); ok || i == len(buffer)-1 {
			sg.mu.RLock()
			messages := messageBuffer{
				bts:      w.Bytes(op),
				messages: w.Messages(op),
			}
			for _, s := range sg.sessions {
				s.Write(messages)
			}
			sg.mu.RUnlock()
			if i != len(buffer)-1 {
				i--
			}
			w.Reset(op)
		}
	}

	sg.mu.Lock()
	sg.intermediateState = false
	if len(sg.deletedSessions) > 0 {
		for _, d := range sg.deletedSessions {
			if s, ok := sg.sessions[d]; ok {
				s.Push(sg.buffer...)
				delete(sg.sessions, d)
			}
		}
		sg.deletedSessions = sg.deletedSessions[:0]
	}
	sg.mu.Unlock()
}

func (sg *subGroup) AddSession(s *session) bool {
	sg.mu.Lock()
	defer sg.mu.Unlock()
	if len(sg.sessions) > sg.maxSize {
		return false
	}
	sg.sessions[s.descriptor] = s
	return true
}

func (sg *subGroup) DeleteSession(d uint64) {
	sg.mu.Lock()
	if sg.intermediateState {
		sg.deletedSessions = append(sg.deletedSessions, d)
	}else if s, exists := sg.sessions[d]; exists {
		s.Push(sg.buffer...)
		delete(sg.sessions, d)
	}
	sg.mu.Unlock()
}

type group struct {
	descriptor uint64
	subs map[uint64]int
	subGroups []*subGroup
	sessionsPerSubGroup int
	publishQueue chan flusher
	mu sync.RWMutex
}

func (g *group) Push(messages ...Message) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, sg := range g.subGroups {
		sg.Push(messages...)
	}
}

func (g *group) Add(s *session) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i := len(g.subGroups) - 1; i >= 0; i-- {
		sg := g.subGroups[i]
		if sg.AddSession(s) {
			g.subs[s.descriptor] = i
			return
		}
	}
	g.subGroups = append(g.subGroups, &subGroup{
		enqueued:          false,
		publishQueue:      g.publishQueue,
		retryEnqueued:     false,
		buffer:            nil,
		sessions: map[uint64]*session{},
		intermediateState: false,
		deletedSessions:   nil,
		maxSize:           1024,
	})
	i := len(g.subGroups) - 1
	g.subGroups[i].AddSession(s)
	g.subs[s.descriptor] = i
}

func (g *group) Len() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.subs)
}

func (g *group) Delete(descriptor uint64) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	fmt.Println("deleting somebody")
	if index, ok := g.subs[descriptor]; ok {
		delete(g.subs, descriptor)
		sg := g.subGroups[index]
		sg.DeleteSession(descriptor)
	} else {
		fmt.Println("JAJAJA LOH")
	}
	return len(g.subs)
}