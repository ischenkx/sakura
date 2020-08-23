package notify

type User struct {
	clients map[string]struct{}
	channels map[string]struct{}
}
