package websockets

import (
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/internal/pubsub"
	"github.com/RomanIschenko/notify/internal/pubsub/internal/client"
	"github.com/RomanIschenko/notify/internal/pubsub/transport"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

var logger = logrus.WithField("source", "ws_logger")

type Server struct {
	acceptChan     chan notify.IncomingConnection
	incomingChan   chan notify.IncomingData
	inactivateChan chan *client.Client
	upgrader       ws.Upgrader
	httpUpgrader   ws.HTTPUpgrader
}

func (s *Server) serveConn(c net.Conn) {


	if tcpConn, ok := c.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
	}


	t := &Transport{
		conn:   c,
		state:  int32(transport.Open),
		uid: uuid.New().String(),
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
			logger.Debugf("(websockets)session.Recv failed, inactivating connection:", err)
			s.inactivateChan <- client
			return
		}

		s.incomingChan <- notify.IncomingData{Client: client, Payload: data}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := s.httpUpgrader.Upgrade(r,w)
	if err != nil {
		logger.Debug(err)
	}

	go s.serveConn(conn)
}

func (s *Server) Accept() <-chan notify.IncomingConnection {
	return s.acceptChan
}

func (s *Server) Inactive() <-chan *client.Client {
	return s.inactivateChan
}

func (s *Server) Incoming() <-chan notify.IncomingData {
	return s.incomingChan
}

func NewServer(upgrader ws.Upgrader, httpUpgrader ws.HTTPUpgrader) *Server {
	s := &Server{
		acceptChan:     make(chan notify.IncomingConnection, 1024),
		incomingChan:   make(chan notify.IncomingData, 1024),
		inactivateChan: make(chan *client.Client, 1024),
		upgrader:       upgrader,
		httpUpgrader:   httpUpgrader,
	}
	return s
}