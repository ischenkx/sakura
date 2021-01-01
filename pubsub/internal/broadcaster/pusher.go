package broadcaster

type pusher interface {
	Push(...Message)
}
