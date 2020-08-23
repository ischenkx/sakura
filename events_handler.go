package notify

import (
	"bufio"
	"io"
	"sync"
)

type JoinHandler func(JoinOptions)
type LeaveHandler func(LeaveOptions)
type SendHandler func(SendOptions)
type DataHandler func(client *Client, reader io.Reader) error

type EventsHandler struct {
	joinHandler  	JoinHandler
	leaveHandler 	LeaveHandler
	sendHandler  	SendHandler
	dataHandler		DataHandler

	sysJoinHandler 	JoinHandler
	sysLeaveHandler LeaveHandler
	sysSendHandler  SendHandler

	mu 				sync.RWMutex
}

func (eventsHandler *EventsHandler) Join(handler JoinHandler) {
	if handler == nil {
		return
	}
	eventsHandler.mu.Lock()
	defer eventsHandler.mu.Unlock()

	eventsHandler.joinHandler = handler
}

func (eventsHandler *EventsHandler) Leave(handler LeaveHandler) {
	if handler == nil {
		return
	}
	eventsHandler.mu.Lock()
	defer eventsHandler.mu.Unlock()

	eventsHandler.leaveHandler = handler
}

func (eventsHandler *EventsHandler) Send(handler SendHandler) {
	if handler == nil {
		return
	}
	eventsHandler.mu.Lock()
	defer eventsHandler.mu.Unlock()

	eventsHandler.sendHandler = handler
}

func (eventsHandler *EventsHandler) Data(handler DataHandler) {
	if handler == nil {
		return
	}
	eventsHandler.mu.Lock()
	defer eventsHandler.mu.Unlock()

	eventsHandler.dataHandler = handler
}

func (eventsHandler *EventsHandler) sysJoin(handler JoinHandler) {
	if handler == nil {
		return
	}
	eventsHandler.mu.Lock()
	defer eventsHandler.mu.Unlock()

	eventsHandler.sysJoinHandler = handler
}

func (eventsHandler *EventsHandler) sysLeave(handler LeaveHandler) {
	if handler == nil {
		return
	}
	eventsHandler.mu.Lock()
	defer eventsHandler.mu.Unlock()

	eventsHandler.sysLeaveHandler = handler
}

func (eventsHandler *EventsHandler) sysSend(handler SendHandler) {
	if handler == nil {
		return
	}
	eventsHandler.mu.Lock()
	defer eventsHandler.mu.Unlock()

	eventsHandler.sysSendHandler = handler
}


func (eventsHandler *EventsHandler) handle(opts interface{}) {
	if opts == nil {
		return
	}

	switch opts.(type) {
	case JoinOptions:
		eventsHandler.mu.RLock()
		h1, h2 := eventsHandler.sysJoinHandler, eventsHandler.joinHandler
		eventsHandler.mu.RUnlock()
		if h1 != nil {
			h1(opts.(JoinOptions))
		}
		if h2 != nil {
			h2(opts.(JoinOptions))
		}
	case LeaveOptions:
		eventsHandler.mu.RLock()
		h1, h2 := eventsHandler.sysLeaveHandler, eventsHandler.leaveHandler
		eventsHandler.mu.RUnlock()
		if h1 != nil {
			h1(opts.(LeaveOptions))
		}
		if h2 != nil {
			h2(opts.(LeaveOptions))
		}
	case SendOptions:
		eventsHandler.mu.RLock()
		h1, h2 := eventsHandler.sysSendHandler, eventsHandler.sendHandler
		eventsHandler.mu.RUnlock()
		if h1 != nil {
			h1(opts.(SendOptions))
		}
		if h2 != nil {
			h2(opts.(SendOptions))
		}
	}
}

func (eventsHandler *EventsHandler) handleData(client *Client, reader io.Reader) error {
	var err error = nil
	eventsHandler.mu.RLock()
	handler := eventsHandler.dataHandler
	eventsHandler.mu.RUnlock()
	if handler != nil {
		err = handler(client, reader)
		return err
	}

	_, err = bufio.NewReader(reader).Peek(1)
	return err
}