package sockjs

import (
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/igm/sockjs-go/v3/sockjs"
)

type Transport struct {
	session sockjs.Session
	client  pubsub.Client
}

func (t *Transport) Write(d []byte) (int, error) {
	return len(d), t.session.Send(string(d))
}

func (t *Transport) Close() error {
	return t.session.Close(0, "")
}

func (t *Transport) Open() bool {
	return t.session.GetSessionState() != sockjs.SessionClosed
}
