package notify

import "time"

type servable App

func (s *servable) Connect(auth string, options ConnectOptions) (Client, error){
	return (*App)(s).connect(options, auth)
}

func (s *servable) Inactivate(id string) {
	client, changelog, err := s.pubsub.Inactivate(id, time.Now().UnixNano())
	if err != nil {
		return
	}
	s.eventsRegistry.emitInactivate((*App)(s), client)
	s.eventsRegistry.emitChange((*App)(s), changelog)
}

func (s *servable) HandleMessage(client Client, data []byte) {
	s.eventsRegistry.emitMessage((*App)(s), client, data)
}



type Servable interface {
	Connect(auth string, options ConnectOptions) (Client, error)
	HandleMessage(client Client, data []byte)
	Inactivate(id string)
}
