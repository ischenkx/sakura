package broadcaster

type flusher interface {
	Flush(w *messageWriter)
}
