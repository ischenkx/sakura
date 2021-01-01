package broadcaster

type Message struct {
	Data []byte
	NoBuffering bool
	retries int
}
