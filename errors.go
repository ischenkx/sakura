package swirl

import "fmt"

//type FilterNotPassedErr struct {
//	Reason error
//	Client Client
//	Event Event
//}
//
//func (e FilterNotPassedErr) Error() string {
//	return fmt.Sprintf("client (%s) failed to pass a filter: %s", e.Client.ID(), e.Reason)
//}

type ConnectionNotEstablishedErr struct {
	Reason         error
	ConnectOptions ConnectOptions
}

type DecodingError struct {
	Reason  error
	Client  Client
	Payload []byte
}

type EncodingError struct {
	Reason       error
	EventOptions EventOptions
}

type HandlerInitializationError struct {
	Reason error
	EventName string
	Handler interface{}
}

type HandlerNotFoundError struct {
	Event string
}

type HandlerCallError struct {
	Reason error
}

type FailedMessageRecoveryError struct {
	Topics []string
	Client Client
	TimeStamp int64
}

func (f FailedMessageRecoveryError) Error() string {
	return fmt.Sprintf("failed to recover data for %s (%d)", f.Client.ID(), f.TimeStamp)
}

func (e ConnectionNotEstablishedErr) Error() string {
	return fmt.Sprintf("failed to establish connection: %s", e.Reason)
}

func (e DecodingError) Error() string {
	return fmt.Sprintf("invalid message received from client (%s): %s", e.Client.ID(), e.Reason)
}

func (e EncodingError) Error() string {
	return fmt.Sprintf("failed to encode a message: %s", e.Reason)
}

func (e HandlerNotFoundError) Error() string {
	return fmt.Sprintf("failed to find a handler for event \"%s\"", e.Event)
}

func (e HandlerCallError) Error() string {
	return fmt.Sprintf("failed to call a handler: %s", e.Reason)
}

func (h HandlerInitializationError) Error() string {
	return fmt.Sprintf("failed to initialize a handler for event \"%s\": %s", h.EventName, h.Reason)
}

