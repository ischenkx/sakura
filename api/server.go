package api

import (
	"context"
	"errors"
	"github.com/RomanIschenko/notify"
	notifier "github.com/RomanIschenko/notify/api/pb"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/RomanIschenko/notify/pubsub/publication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var NoBrokerProvidedErr = errors.New("no broker provided")

var empty = new(notifier.EmptyMessage)

var logger = logrus.WithField("source", "notify_api")

type Config struct {
	APIKey string
	Auth   notify.Auth
	Broker notify.Broker
}

func (cfg Config) validate() Config {
	return cfg
}

// grpc server to send manage your pubsub
type Server struct {
	auth notify.Auth
	apiKey string
	broker notify.Broker
	notifier.UnimplementedNotifyServer
}

func (s *Server) Publish(ctx context.Context, req *notifier.PublishRequest) (*notifier.EmptyMessage, error) {
	if s.broker == nil {
		return empty, NoBrokerProvidedErr
	}

	opts := pubsub.PublishOptions{
		Topics:  req.Topics,
		Clients: req.Clients,
		Users:   req.Users,
		Payload: publication.New(req.Message),
		Time:    req.Time,
	}

	s.broker.Emit(notify.BrokerEvent{
		Data:     opts,
		AppID:    req.AppID,
		Event: 	  notify.PublishEvent,
		Time:     req.Time,
	})

	return empty, nil
}

func (s *Server) Unsubscribe(ctx context.Context, req *notifier.UnsubscribeRequest) (*notifier.EmptyMessage, error) {
	if s.broker == nil {
		return nil, NoBrokerProvidedErr
	}

	opts := pubsub.UnsubscribeOptions{
		Topics:  req.Topics,
		Clients: req.Clients,
		Users:   req.Users,
		Time:    req.Time,
		All: 	 req.All,
	}

	s.broker.Emit(notify.BrokerEvent{
		Data:     opts,
		AppID:    req.AppID,
		Event: 	  notify.UnsubscribeEvent,
		Time:     req.Time,
	})

	return empty, nil

}

func (s *Server) Subscribe(ctx context.Context, req *notifier.SubscribeRequest) (*notifier.EmptyMessage, error) {
	if s.broker == nil {
		return nil, NoBrokerProvidedErr
	}
	opts := pubsub.SubscribeOptions{
		Topics:  req.Topics,
		Clients: req.Clients,
		Users:   req.Users,
		Time:    req.Time,
	}

	s.broker.Emit(notify.BrokerEvent{
		Data:     opts,
		AppID:    req.AppID,
		Event: 	  notify.SubscribeEvent,
		Time:     req.Time,
	})

	return empty, nil
}

func (s *Server) Authorize(ctx context.Context, input *notifier.AuthInput) (*notifier.AuthOutput, error) {
	authData := input.ClientID
	if s.auth != nil {
		var err error
		authData, err = s.auth.Register(pubsub.ClientID(input.ClientID))
		if err != nil {
			logrus.Debug("failed to register a client:", err)
			return nil, err
		}
	}
	return &notifier.AuthOutput{
		Token: authData,
	}, nil
}

func (s *Server) Start(ctx context.Context, lis net.Listener, opts... grpc.ServerOption) {
	grpcServer := grpc.NewServer(opts...)

	notifier.RegisterNotifyServer(grpcServer, s)

	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			logger.Errorf("grpc server failed to start:", err)
		}
	}()

	select {
	case <-ctx.Done():
	}
	grpcServer.GracefulStop()
}

// creates new grpc server
func NewServer(config Config) *Server {
	config = config.validate()

	return &Server{
		broker: config.Broker,
		auth:   config.Auth,
		apiKey: config.APIKey,
	}
}