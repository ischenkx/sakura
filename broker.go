package notify

import (
	"time"
)

type BrokerEvent int

const (
	SendEvent BrokerEvent = iota + 1
	JoinEvent
	LeaveEvent
	ConnectEvent
	DisconnectEvent
	InstanceUpEvent
	InstanceDownEvent
	AppUpEvent
	AppDownEvent
)

type BrokerAction struct {
	Data interface{}
	AppID string
	BrokerID string
	Event BrokerEvent
	Time time.Time
}

type BrokerEventHandler func(BrokerAction) error

type Broker interface {
	// Handle register a new handler for broker events
	// if handler returns non-nil error then handler is deleted
	Handle(...BrokerEventHandler)
	Emit(BrokerAction) error
	ID() string
}