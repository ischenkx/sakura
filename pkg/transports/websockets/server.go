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
	app      notify.Servable
	upgrader ws.HTTPUpgrader
}

func (s *Server) authenticate(t *Transport) (notify.Client, error) {
	mes, err := t.readMessage()

	if err != nil {
		return nil, err
	}

	switch mes.opCode {
	case authReqCode:
		authData := string(mes.data)
		return s.app.Connect(authData, notify.ConnectOptions{
			Writer:    t,
			TimeStamp: time.Now().UnixNano(),
			Meta:      t.conn,
		})
	case pingCode:
		err = s.sendPong(mes.data, t)
		if err != nil {
			return nil, err
		}
		return s.authenticate(t)
	default:
		return nil, errors.New("failed to authenticate: received message with a wrong opcode")
	}

}

func (s *Server) handleMessage(c notify.Client, t *Transport, m message) error {
	var err error

	switch m.opCode {
	case pingCode:
		err = s.sendPong(m.data, t)
	case messageCode:
		s.app.HandleMessage(c, m.data)
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
	conn, _, _, err := s.upgrader.Upgrade(r, w)

	if err != nil {
		logger.Println("failed to upgrade connection")
		return
	}

	t := newTransport(conn)

	client, err := s.authenticate(t)

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
		s.app.Inactivate(client.ID())
		t.Close()
		return
	}

	for {
		mes, err := t.readMessage()
		if err != nil {
			s.app.Inactivate(client.ID())
			t.Close()
			return
		}

		err = s.handleMessage(client, t, mes)

		if err != nil {
			s.app.Inactivate(client.ID())
			t.Close()
			return
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveWS(r, w)
}

func NewServer(app notify.Servable, upgrader ws.HTTPUpgrader) *Server {
	s := &Server{}
	s.app = app
	s.upgrader = upgrader
	return s
}
