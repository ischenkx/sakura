package pubsub

import (
	"errors"
	"github.com/RomanIschenko/notify/internal/batch_queue"
	"github.com/RomanIschenko/notify/pubsub/client_id"
	"github.com/RomanIschenko/notify/pubsub/transport"
	"sync"
)

type ClientState int

const (
	Active ClientState = iota
	Inactive
	Invalid
)

type Client struct {
	id             clientid.ID
	queue          *batchqueue.Queue
	hash, userHash int
	state          ClientState
	meta           *sync.Map
	mu             sync.Mutex
}

func (c *Client) tryActivate(t transport.Transport) error {
	if t == nil {
		return errors.New("nil transport activation")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == Invalid {
		return errors.New("can't activate invalid client")
	}

	if t.State() == transport.Closed {
		return errors.New("cannot activate client with closed transport")
	}

	if c.state == Active {
		c.queue.Inactivate()
	}

	c.queue.Activate(t)
	c.state = Active

	return nil
}

func (c *Client) tryInactivate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state != Active {
		return errors.New("can't inactivate non-active client")
	}
	c.state = Inactive
	c.queue.Inactivate()
	return nil
}

func (c *Client) tryInvalidate() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == Invalid {
		return nil
	}
	if c.state == Active {
		return errors.New("can't invalidate client with active state, actually you can force invalidate")
	}
	c.queue.Inactivate()
	c.state = Invalid
	return nil
}

func (c *Client) forceInvalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == Invalid {
		return
	}
	c.queue.Inactivate()
	c.state = Invalid
}

func (c *Client) publish(p []byte) {
	c.queue.Push(p)
}

func (c *Client) ID() clientid.ID {
	return c.id
}

func (c *Client) Hash() int {
	return c.hash
}

func (c *Client) UserHash() int {
	return c.userHash
}

func (c *Client) Meta() *sync.Map {
	return c.meta
}

func newClient(id clientid.ID, bufferSize int) (*Client, error) {
	idHash, err := id.Hash()

	if err != nil {
		return nil, err
	}

	return &Client{
		id:    id,
		hash:  idHash,
		meta:  &sync.Map{},
		queue: batchqueue.New(bufferSize),
		state: Inactive,
	}, nil
}