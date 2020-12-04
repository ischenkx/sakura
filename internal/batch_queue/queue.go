package batchqueue

import (
	"bytes"
	"encoding/binary"
	"io"
	"sync"
)

var pubsPool = &sync.Pool{
	New: func() interface{} {
		return make([][]byte, 256)
	},
}

type Queue struct {
	mu            sync.RWMutex
	writer		  io.Writer
	pubs          [][]byte
	maxBufferSize int
	active        bool
	sig           *sync.Cond
}

func (q *Queue) Activate(w io.Writer) {
	q.mu.Lock()
	q.writer = w
	q.active = true
	q.mu.Unlock()
	go q.run()
}

func (q *Queue) Inactivate() {
	q.mu.Lock()
	if q.active {
		q.active = false
		q.sig.Signal()
	}
	q.mu.Unlock()
}

func (q *Queue) Push(data ...[]byte) {
	if len(data) == 0 {
		return
	}
	q.mu.Lock()
	q.pubs = append(q.pubs, data...)
	q.sig.Signal()
	q.mu.Unlock()
}

func (q *Queue) run() {
	buffer := bufferPool.Get().(*bytes.Buffer)
	var failedPubs [][]byte
	var largePubs [][]byte
	defer bufferPool.Put(buffer)
	q.mu.Lock()
	maxSize := q.maxBufferSize
	batch := newBatch(buffer, q.maxBufferSize)
	q.mu.Unlock()

	for {
		q.mu.Lock()
		if len(q.pubs) == 0 {
			q.sig.Wait()
		}
		if !q.active {
			q.mu.Unlock()
			return
		}

		pubs := q.pubs
		q.pubs = pubsPool.Get().([][]byte)
		w := q.writer
		q.mu.Unlock()
		for i := 0; i < len(pubs); i++ {
			pub := pubs[i]
			if len(pub) >= maxSize {
				largePubs = append(largePubs, pub)
				continue
			}
			if successful := batch.Push(pub); !successful {
				_, err := w.Write(batch.Bytes())
				if err != nil {
					failedPubs = append(failedPubs, batch.Fallback()...)
				}
				batch.Reset()
				i--
			}
		}
		if batch.curSize > 0 {
			_, err := w.Write(batch.Bytes())
			if err != nil {
				failedPubs = append(failedPubs, batch.Fallback()...)
			}
			batch.Reset()
		}
		if len(failedPubs) > 0 {
			q.Push(failedPubs...)
			failedPubs = failedPubs[:0]
		}
		if len(largePubs) > 0 {
			for _, pub := range largePubs {
				binary.Write(buffer, binary.LittleEndian, uint32(len(pub)))
				buffer.Write(pub)
				w.Write(buffer.Bytes())
				buffer.Reset()
			}
			largePubs = largePubs[:0]
		}
		pubsPool.Put(pubs[:0])
	}
}

func New(maxBufferSize int) *Queue {
	if maxBufferSize <= 0 {
		maxBufferSize = 256
	}

	queue := &Queue{
		mu:            sync.RWMutex{},
		pubs:          pubsPool.Get().([][]byte),
		maxBufferSize: maxBufferSize,
	}

	queue.sig = sync.NewCond(&queue.mu)

	return queue
}