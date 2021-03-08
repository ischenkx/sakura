package websockets

import (
	"errors"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
)

type Transport struct {
	conn net.Conn
}

func (t *Transport) Write(d []byte) (int, error) {
	if t.conn == nil {
		return 0, errors.New("connection is closed or not open")
	}
	return len(d), wsutil.WriteServerMessage(t.conn, ws.OpText, d)
}

func (t *Transport) Close() error {
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

func (t *Transport) Open() bool {
	return t.conn != nil
}
