package notify

type Auth interface {
	Authorize(token string) (id string, err error)
	Register(id string) (token string, err error)
}
