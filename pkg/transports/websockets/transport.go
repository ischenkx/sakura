package websockets

import (
	"bytes"
	"errors"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

type Transport struct {
	conn net.Conn
	buf  *bytes.Buffer
	mu   sync.Mutex
}

// writes a message
func (t *Transport) Write(d []byte) (int, error) {
	return t.writeServiceMessage(message{
		opCode: messageCode,
		data:   d,
	})
}

func (t *Transport) writeServiceMessage(m message) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return 0, errors.New("already closed")
	}

	if _, err := writeMessage(t.buf, m); err != nil {
		t.buf.Reset()
		return 0, err
	}
	bts := t.buf.Bytes()
	n := len(bts)
	err := wsutil.WriteServerBinary(t.conn, bts)
	if m.opCode == authAckCode && len(m.data) > 8 {
		log.Println("sending:", string(m.data))
	}
	t.buf.Reset()
	return n, err
}

func (t *Transport) readMessage() (message, error) {
	if !t.Open() {
		return message{}, errors.New("transport is already closed")
	}
	bts, _, err := wsutil.ReadClientData(t.conn)
	if err != nil {
		return message{}, err
	}
	return readMessage(bts)
}

func (t *Transport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.conn == nil {
		return nil
	}

	t.buf.Reset()
	bufferPool.Put(t.buf)

	err := t.conn.Close()

	t.buf = nil
	t.conn = nil

	return err
}

func (t *Transport) Open() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.conn != nil
}

func newTransport(conn net.Conn) *Transport {
	return &Transport{
		conn: conn,
		buf:  bufferPool.Get().(*bytes.Buffer),
	}
}
