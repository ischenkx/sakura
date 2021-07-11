package swirl

import (
	"io"
	"log"
	"time"
)

type ConnectOptions struct {
	Auth      string
	Writer    io.WriteCloser
	TimeStamp int64
	Meta      interface{}
}

type Server struct {
	app *App
}

func (s *Server) Inactivate(id string, ts int64) {
	s.app.pubsub.Inactivate(id, ts)
}

func (s *Server) Connect(opts ConnectOptions) (Client, error) {
	cl := ChangeLog{TimeStamp: opts.TimeStamp}
	clientID, userID, err := s.app.auth.Authorize(opts.Auth)

	if err != nil {
		s.app.events.callError(ConnectionNotEstablishedErr{
			Reason:         err,
			ConnectOptions: opts,
		})
		return nil, err
	}
	res, err := s.app.pubsub.Connect(clientID, opts.Writer, opts.TimeStamp)
	if err != nil {
		s.app.events.callError(ConnectionNotEstablishedErr{
			Reason:         err,
			ConnectOptions: opts,
		})
		return nil, err
	}

	cl.ClientsUp = append(cl.ClientsUp, clientID)

	cl.Merge(s.app.pubsub.Users().Add(userID, clientID, opts.TimeStamp))

	for _, sub := range s.app.User(clientID).Subscriptions().Array() {
		if c, err := s.app.pubsub.Subscribe(clientID, sub.Topic, sub.TimeStamp); err == nil {
			cl.Merge(c)
		}
	}

	s.app.events.callChange(cl)

	if res.IsReconnect {
		s.app.events.callReconnect(opts, s.app.Client(clientID))
	} else {
		s.app.events.callConnect(opts, s.app.Client(clientID))
	}

	if len(res.UnrecoverableTopics) > 0 {
		s.app.events.callError(FailedMessageRecoveryError{
			Topics:    res.UnrecoverableTopics,
			Client:    s.app.Client(clientID),
			TimeStamp: opts.TimeStamp,
		})
	}

	return s.app.Client(clientID), nil
}

func (s *Server) HandleMessage(clientID string, message []byte) {

	name, argBytes, err := s.app.emitter.DecodeRawData(message)

	if err != nil {
		log.Println(err)
		return
	}

	h, ok := s.app.emitter.GetHandler(name)

	if !ok {
		s.app.events.callError(HandlerNotFoundError{name})
		return
	}

	args, err := h.DecodeData(argBytes)

	if err != nil {
		s.app.events.callError(DecodingError{
			Reason:  err,
			Client:  s.app.Client(clientID),
			Payload: message,
		})
		return
	}

	s.app.events.callEvent(s.app.Client(clientID), EventOptions{
		Name:      name,
		Args:      args,
		TimeStamp: time.Now().UnixNano(),
	})

	err = h.Call(s.app, s.app.Client(clientID), args...)
	if err != nil {
		s.app.events.callError(HandlerCallError{Reason: err})
		return
	}
}
