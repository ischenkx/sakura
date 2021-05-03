package notify

type servable App

func (s *servable) Connect(auth string, options ConnectOptions) (Client, error) {
	return (*App)(s).connect(options, auth)
}

func (s *servable) Inactivate(id string) {
	(*App)(s).inactivate(id)
}

func (s *servable) HandleMessage(client Client, data []byte) {
	(*App)(s).handleMessage(data, client)
}

type Servable interface {
	Connect(auth string, options ConnectOptions) (Client, error)
	HandleMessage(client Client, data []byte)
	Inactivate(id string)
}
