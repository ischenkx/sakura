package notify

type server App

func (s *server) Connect(auth string, options ConnectOptions) (string, error) {
	return (*App)(s).connect(options, auth)
}

func (s *server) Inactivate(id string) {
	(*App)(s).inactivate(id)
}

func (s *server) HandleMessage(id string, data []byte) {
	(*App)(s).handleMessage(id, data)
}

type Server interface {
	Connect(auth string, options ConnectOptions) (string, error)
	HandleMessage(id string, data []byte)
	Inactivate(id string)
}
