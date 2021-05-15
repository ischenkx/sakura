package websockets

import (
	"errors"
	"github.com/RomanIschenko/notify"
	"github.com/gobwas/ws"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"sync"
	"time"
)

var logger = logrus.WithField("source", "ws_logger")

type Config struct {
	AuthReadTimeout time.Duration
	ReadTimeout time.Duration
	Upgrader ws.HTTPUpgrader
}

// Server uses a custom protocol for websocket communication.
//
// Protocol:
// (opcode (1 byte))(data (n bytes))
//
// Connection algorithm is pretty simple:
//
// 1. Client connects
//
// 2. Client sends auth data
//
// 3. Server authenticates client
//
// 4.
//    Authentication failure => Server sends a message with a specified opcode and data - "fail"
//    Authentication succeeds => Server sends a message with a specified opcode and data - "ok"
//
// 5. Server reads messages and handles ping requests
type Server struct {
	mu       sync.RWMutex
	app      notify.Server
	config Config
}

func (s *Server) authenticate(t *Transport) (string, error) {

	deadline := time.Time{}

	if s.config.AuthReadTimeout > 0 {
		deadline = time.Now().Add(s.config.AuthReadTimeout)
	}

	if err := t.conn.SetReadDeadline(deadline); err != nil {
		return "", err
	}

	mes, err := t.readMessage()

	if err != nil {
		return "", err
	}

	switch mes.opCode {
	case authReqCode:
		authData := string(mes.data)
		t.conn.SetReadDeadline(time.Time{})
		return s.app.Connect(authData, notify.ConnectOptions{
			Writer:    t,
			TimeStamp: time.Now().UnixNano(),
			Meta:      t.conn,
		})
	case pingCode:
		err = s.sendPong(mes.data, t)
		if err != nil {
			return "", err
		}
		return s.authenticate(t)
	default:
		return "", errors.New("failed to authenticate: received message with a wrong opcode")
	}

}

func (s *Server) handleMessage(id string, t *Transport, m message) error {
	var err error

	switch m.opCode {
	case pingCode:
		err = s.sendPong(m.data, t)
	case messageCode:
		s.app.HandleMessage(id, m.data)
	default:
		log.Println("unexpected code:", m.opCode)
	}

	return err
}

func (s *Server) sendPong(data []byte, t *Transport) error {
	_, err := t.writeServiceMessage(message{
		opCode: pongCode,
		data:   data,
	})

	return err
}

func (s *Server) serveWS(r *http.Request, w http.ResponseWriter) {
	conn, _, _, err := s.config.Upgrader.Upgrade(r, w)

	if err != nil {
		logger.Println("failed to upgrade connection")
		return
	}

	t := newTransport(conn)

	id, err := s.authenticate(t)

	if err != nil {
		log.Println(err)
		t.writeServiceMessage(message{
			opCode: authAckCode,
			data:   []byte("fail"),
		})
		t.Close()
		return
	}

	_, err = t.writeServiceMessage(message{
		opCode: authAckCode,
		data:   []byte("ok"),
	})

	if err != nil {
		s.app.Inactivate(id)
		t.Close()
		return
	}

	for {
		deadline := time.Time{}
		if s.config.ReadTimeout > 0 {
			deadline = time.Now().Add(s.config.ReadTimeout)
		}
		err = t.conn.SetReadDeadline(deadline)
		if err != nil {
			s.app.Inactivate(id)
			t.Close()
			return
		}
		mes, err := t.readMessage()
		if err != nil {
			s.app.Inactivate(id)
			t.Close()
			return
		}
		err = s.handleMessage(id, t, mes)
		if err != nil {
			s.app.Inactivate(id)
			t.Close()
			return
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveWS(r, w)
}

func NewServer(app notify.Server, config Config) *Server {
	s := &Server{}
	s.app = app
	s.config = config
	return s
}
