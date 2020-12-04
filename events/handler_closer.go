package events

import (
	"errors"
)

type HandlerCloser struct {
	uid, event string
	src *Source
}

func (c HandlerCloser) Close() error {
	if c.src == nil {
		return errors.New("failed to close handler")
	}
	c.src.mu.Lock()
	defer c.src.mu.Unlock()
	delete(c.src.handlers[c.event], c.uid)
	return nil
}
