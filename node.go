package notify

import (
	"io"
	"notify/cleaner"
	"sync"
	"time"
)

const (
	RunningNode int32 = iota + 1
	StoppedNode
)

type NodeConfig struct {
	ID string
	Messages MessagesStorage
	Broker Broker
}

type Node struct {
	apps     map[string]*App
	mu       sync.RWMutex
	id       string
	messages MessagesStorage
	broker   Broker
	cleaner  *cleaner.Cleaner
	state    int32
}

func (node *Node) Run() {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.state == RunningNode {
		return
	}
	node.cleaner.Run()
	for _, app := range node.apps {
		app.Run()
	}
	node.state = RunningNode
	node.broker.Handle(func(mes BrokerMessage) {
		if node == nil {
			return
		}
		app := node.GetApp(mes.AppID)
		if app == nil {
			return
		}
		switch mes.Event {
		case Send:
			if opts, ok := mes.Data.(SendOptions); ok {
				app.send(opts)
			}
		case Join:
			if opts, ok := mes.Data.(JoinOptions); ok {
				app.join(opts)
			}
		case Leave:
			if opts, ok := mes.Data.(LeaveOptions); ok {
				app.leave(opts)
			}
		}
	})
}

func (node *Node) Stop() {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.state == StoppedNode {
		return
	}
	node.cleaner.Stop()
	for _, app := range node.apps {
		app.Stop()
	}
	node.state = StoppedNode

}


func (node *Node) GetApp(id string) *App {
	node.mu.RLock()
	app := node.apps[id]
	node.mu.RUnlock()
	return app
}

func (node *Node) handleAppEvents(app *App) {
	if app == nil {
		return
	}
	app.On().sysJoin(func(opts JoinOptions) {
		if node == nil {
			return
		}
		if node.broker == nil {
			return
		}
		node.broker.Send(BrokerMessage{
			Data:   opts,
			AppID:  app.id,
			Event:  Join,
		})
	})
	app.On().sysLeave(func(opts LeaveOptions) {
		if node == nil || app == nil {
			return
		}
		if node.broker == nil {
			return
		}
		node.broker.Send(BrokerMessage{
			Data:   opts,
			AppID:  app.id,
			Event:  Leave,
		})
	})
	app.On().sysSend(func(opts SendOptions) {
		if node == nil {
			return
		}
		if node.broker == nil {
			return
		}
		node.broker.Send(BrokerMessage{
			Data:   opts,
			AppID:  app.id,
			Event:  Send,
		})
	})
	app.On().Data(func(client *Client, reader io.Reader) error {
		
		return nil
	})

}

func (node *Node) App(id string) *App {
	node.mu.RLock()
	app, ok := node.apps[id]
	node.mu.RUnlock()

	if !ok {
		app = NewApp(AppConfig{
			ID:       id,
			Messages: node.messages,
			Broker: node.broker,
			controlled: true,
		})
		app.cleaner = node.cleaner
		node.mu.Lock()
		a, existing := node.apps[id]
		if !existing {
			node.apps[id] = app
		} else {
			app = a
		}
		node.mu.Unlock()
		if !existing {
			node.handleAppEvents(app)
			app.Run()
		}
	}
	return app
}

func (node *Node) Clean() {
	node.mu.RLock()
	apps := make([]*App, 0, len(node.apps))
	for _, app := range node.apps {
		apps = append(apps, app)
	}
	node.mu.RUnlock()

	for _, app := range apps {
		app.Clean()
	}
}

func NewNode(config NodeConfig) *Node {
	node := &Node{
		apps: 	  map[string]*App{},
		mu:       sync.RWMutex{},
		id:       config.ID,
		messages: config.Messages,
		broker:   config.Broker,
		state:    StoppedNode,
	}

	node.cleaner = cleaner.New(node, time.Minute * 4)

	return node
}