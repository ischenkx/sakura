package batchproto

import (
	"github.com/ischenkx/swirl/internal/pubsub/message"
	"github.com/ischenkx/swirl/internal/pubsub/protocol"
)

type Protocol struct {
	batcher *batcher
}

func (p Protocol) Encode(w protocol.Writer, b message.Buffer, meta interface{}) {
	p.batcher.PutMessages(b)
	for {
		data, ok := p.batcher.Next()
		if !ok {
			break
		}
		w.Write(message.NewBatch(data.Messages(), data.Bytes()), meta)
	}
	p.batcher.Reset()
}
