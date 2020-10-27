package pubsub

import (
	"errors"
	"github.com/RomanIschenko/notify/pubsub/publication"
	"sync"
)

type ClientState int

const (
	Active   ClientState = iota
	Inactive
	Invalid
)

type Client struct {
	id		  ClientID
	transport Transport
	buffer publication.Buffer
	hash, userHash int
	state 	   ClientState
	mu 		   sync.Mutex
}

func (c *Client) tryActivate(t Transport) error {
	if t == nil {
		return errors.New("nil transport activation")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == Invalid {
		return errors.New("can't activate invalid client")
	}
	if t.State() == ClosedTransport {
		return errors.New("cannot activate client with closed transport")
	}
	if c.state == Active {
		c.transport.Close()
	}
	c.state = Active
	c.transport = t
	return nil
}

func (c *Client) tryInactivate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state != Active {
		return errors.New("can't inactivate non-active client")
	}
	c.state = Inactive
	c.transport.Close()
	c.transport = nil
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
	if c.transport != nil {
		c.transport.Close()
		c.transport = nil
	}
	c.state = Invalid
	return nil
}

func (c *Client) forceInvalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == Invalid {
		return
	}
	if c.transport != nil {
		c.transport.Close()
		c.transport = nil
	}
	c.state = Invalid
}

func (c *Client) Publish(p publication.Publication) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch c.state {
	case Invalid:
		return errors.New("can't publish to invalid client")
	case Inactive:
		c.buffer.Push(p)
	case Active:
		_, err := c.transport.Write(p.Data)
		return err
	}
	return nil
}

func (c *Client) ID() ClientID {
	return c.id
}

func (c *Client) Hash() int {
	return c.hash
}

func (c *Client) UserHash() int {
	return c.userHash
}

func newClient(id ClientID, bufferSize int) (*Client, error) {
	idHash, err := id.Hash()

	if err != nil {
		return nil, err
	}

	return &Client{
		id:  	id,
		hash:   idHash,
		buffer: publication.NewBuffer(bufferSize),
		state:  Inactive,
	}, nil
}