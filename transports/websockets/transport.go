package websockets

import (
	"github.com/RomanIschenko/notify/internal/pubsub/transport"
	"github.com/gobwas/ws/wsutil"
	"net"
	"sync/atomic"
)

type Transport struct {
	conn net.Conn
	state int32
	uid string
}

func (t *Transport) Write(d []byte) (int, error) {
	return len(d), wsutil.WriteServerBinary(t.conn, d)
}

func (t *Transport) Close() error {
	atomic.StoreInt32(&t.state, int32(transport.Closed))
	return t.conn.Close()
}

func (t *Transport) State() transport.State {
	return transport.State(atomic.LoadInt32(&t.state))
}

