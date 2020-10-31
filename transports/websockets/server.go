package websockets

import (
	"bytes"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

var logger = logrus.WithField("source", "ws_logger")

type Server struct {
	acceptChan chan notify.IncomingConnection
	incomingChan chan notify.IncomingData
	inactivateChan chan *pubsub.Client
	upgrader ws.Upgrader
	httpUpgrder ws.HTTPUpgrader

}

func (s *Server) serveConn(c net.Conn) {
	t := &Transport{
		conn:   c,
		state:  int32(pubsub.OpenTransport),
	}
	data, _, err := wsutil.ReadClientData(c)
	if err != nil {
		logger.Errorf("failed to recv data, closing:", err)
		t.Close()
		return
	}

	conn := notify.IncomingConnection{
		Opts: pubsub.ConnectOptions{
			Transport: t,
			Time:      time.Now().UnixNano(),
		},
		AuthData: string(data),
		Resolver: make(chan notify.ResolvedConnection),
	}

	s.acceptChan <- conn
	resolvedConn := <-conn.Resolver

	if resolvedConn.Err != nil {
		logger.Errorf("failed to resolve connection, closing:", resolvedConn.Err)
		t.Close()
		return
	}

	client := resolvedConn.Client

	for {
		data, _, err := wsutil.ReadClientData(c)

		if err != nil {
			logger.Debugf("(sockjs)session.Recv failed, inactivating connection:", err)
			s.inactivateChan <- client
			return
		}

		s.incomingChan <- notify.IncomingData{client, bytes.NewReader(data)}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := s.httpUpgrder.Upgrade(r,w)
	if err != nil {
		logger.Debug(err)
	}

	go s.serveConn(conn)
}

func (s *Server) Accept() <-chan notify.IncomingConnection {
	return s.acceptChan
}

func (s *Server) Inactive() <-chan *pubsub.Client {
	return s.inactivateChan
}

func (s *Server) Incoming() <-chan notify.IncomingData {
	return s.incomingChan
}

func NewServer(upgrader ws.Upgrader, httpUpgrader ws.HTTPUpgrader) *Server {
	s := &Server{
		acceptChan:   make(chan notify.IncomingConnection, 1024),
		incomingChan: make(chan notify.IncomingData, 1024),
		inactivateChan: make(chan *pubsub.Client, 1024),
		upgrader: upgrader,
		httpUpgrder: httpUpgrader,
	}
	return s
}