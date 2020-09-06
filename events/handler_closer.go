package events

type HandlerCloser struct {
	uid string
	src *Source
}

func (c HandlerCloser) Close() {
	if c.src == nil {
		return
	}
	c.src.mu.Lock()
	defer c.src.mu.Unlock()
	delete(c.src.handlers, c.uid)
}
