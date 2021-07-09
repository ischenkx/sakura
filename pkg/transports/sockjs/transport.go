package sockjs

import (
	"github.com/igm/sockjs-go/v3/sockjs"
	"github.com/ischenkx/swirl"
)

type Transport struct {
	session sockjs.Session
	client  swirl.Client
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
