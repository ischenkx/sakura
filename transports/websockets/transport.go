package websockets

import (
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/gobwas/ws/wsutil"
	"net"
	"sync/atomic"
)

var WritesCounter int32 = 0

type Transport struct {
	conn net.Conn
	state int32
}

func (t *Transport) Write(d []byte) (int, error) {
	atomic.AddInt32(&WritesCounter, 1)
	return len(d), wsutil.WriteServerBinary(t.conn, d)
}

func (t *Transport) Close() error {
	atomic.StoreInt32(&t.state, int32(pubsub.ClosedTransport))
	return t.conn.Close()
}

func (t *Transport) State() pubsub.TransportState {
	return pubsub.TransportState(atomic.LoadInt32(&t.state))
}

