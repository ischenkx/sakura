package sockjs

import (
	pubsub2 "github.com/RomanIschenko/notify/default/pubsub"
	"github.com/igm/sockjs-go/v3/sockjs"
)

type Transport struct {
	session sockjs.Session
	client  pubsub2.Client
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
