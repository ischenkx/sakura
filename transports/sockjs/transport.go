package sockjs

import (
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/transport"
	"github.com/igm/sockjs-go/sockjs"
	"sync/atomic"
)

type Transport struct {
	session sockjs.Session
	state	int32
	client *pubsub.Client
}

func (t *Transport) Write(d []byte) (int, error) {
	return len(d), t.session.Send(string(d))
}

func (t *Transport) Close() error {
	atomic.StoreInt32(&t.state, int32(transport.Closed))
	return t.session.Close(0, "")
}

func (t *Transport) State() transport.State {
	return transport.State(atomic.LoadInt32(&t.state))
}



