package batchproto

import (
	message2 "github.com/ischenkx/notify/internal/pubsub/message"
	protocol2 "github.com/ischenkx/notify/internal/pubsub/protocol"
)

type Protocol struct {
	batcher *batcher
}

func (p Protocol) Encode(w protocol2.Writer, b message2.Buffer, meta interface{}) {
	p.batcher.PutMessages(b)
	for {
		data, ok := p.batcher.Next()
		if !ok {
			break
		}
		w.Write(message2.NewBatch(data.Messages(), data.Bytes()), meta)
	}
	p.batcher.Reset()
}
