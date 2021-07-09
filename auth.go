package swirl

type Authorizer interface {
	Authorize(data string) (client, user string, err error)
}
