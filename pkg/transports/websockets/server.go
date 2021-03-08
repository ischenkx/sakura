package websockets

import (
	"github.com/RomanIschenko/notify"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"sync"
	"time"
)

var logger = logrus.WithField("source", "ws_logger")

type Server struct {
	mu 			 sync.RWMutex
	app			 notify.Servable
	upgrader	 ws.HTTPUpgrader
}

func (s *Server) serveWS(r *http.Request, w http.ResponseWriter) {
	conn, _, _, err := s.upgrader.Upgrade(r, w)

	if err != nil {
		logger.Println("failed to upgrade connection")
		return
	}

	t := &Transport{
		conn:  conn,
	}
	auth, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		logger.Errorf("failed to recv data, closing:", err)
		t.Close()
		return
	}

	client, err := s.app.Connect(string(auth), notify.ConnectOptions{
		Writer:    t,
		TimeStamp: time.Now().UnixNano(),
		Meta:      r,
	})

	if err != nil {
		return
	}

	for {
		message, _, err := wsutil.ReadClientData(t.conn)
		if err != nil {
			log.Println("ws connection closed:", err)
			t.Close()
			s.app.Inactivate(client.ID())
			return
		}
		s.app.HandleMessage(client, message)
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