package authmock

type Auth struct{}

func (a Auth) Authorize(id string) (string, string, error) {
	return id, id, nil
}

func (a Auth) Register(id, user string) (string, error) {
	return "", nil
}

func New() Auth {
	return Auth{}
}
