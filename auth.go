package notify

type Auth interface {
	Authorize(token string) (id, userId string, err error)
	Register(id, userId string) (token string, err error)
}
