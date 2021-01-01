package notify

import (
	"context"
	"github.com/RomanIschenko/notify/pubsub"
	"time"
)

type servable App

func (s *servable) Connect(auth string, options pubsub.ConnectOptions) (pubsub.Client, error){
	return (*App)(s).connect(options, auth)
}

func (s *servable) HandleIncomingData(data IncomingData) {
	s.eventsRegistry.emitIncomingData((*App)(s), data)
}

func (s *servable) Inactivate(id string) {
	s.hub.Inactivate(pubsub.InactivateOptions{
		Time:     time.Now().UnixNano(),
		ClientID: id,
	})
}

type ServableApp interface {
	Connect(auth string, options pubsub.ConnectOptions) (pubsub.Client, error)
	HandleIncomingData(data IncomingData)
	Inactivate(id string)
}

type IncomingData struct {
	Client  pubsub.Client
	Payload []byte
}

type Server interface {
	Start(ctx context.Context, app ServableApp)
}
