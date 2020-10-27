package events

import "errors"

type HandlerCloser struct {
	uid string
	src *Source
}

func (c HandlerCloser) Close() error {
	if c.src == nil {
		return errors.New("failed to close handler")
	}
	c.src.mu.Lock()
	defer c.src.mu.Unlock()
	if h, ok := c.src.handlers[c.uid]; ok {
		delete(c.src.handlers, c.uid)
		close(h.closer)
	}
	return nil
}
