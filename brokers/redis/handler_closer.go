package redibroker

type HandlerCloser struct {
	broker *Broker
	uid string
}

func (c HandlerCloser) Close() error {
	if c.broker != nil {
		c.broker.mu.Lock()
		defer c.broker.mu.Unlock()
		delete(c.broker.handlers, c.uid)
	}
	return nil
}
