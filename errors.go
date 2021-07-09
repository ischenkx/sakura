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

func (e ConnectionNotEstablishedErr) Error() string {
	return fmt.Sprintf("failed to establish connection: %s", e.Reason)
}

type DecodingError struct {
	Reason  error
	Client  Client
	Payload []byte
}

func (e DecodingError) Error() string {
	return fmt.Sprintf("invalid message received from client (%s): %s", e.Client.ID(), e.Reason)
}

type EncodingError struct {
	Reason       error
	EventOptions EventOptions
}

func (e EncodingError) Error() string {
	return fmt.Sprintf("failed to encode a message: %s", e.Reason)
}

type HandlerNotFoundErr struct {
	Event string
}

func (e HandlerNotFoundErr) Error() string {
	return fmt.Sprintf("failed to find a handler for event \"%s\"", e.Event)
}

type HandlerCallError struct {
	Reason error
}

func (e HandlerCallError) Error() string {
	return fmt.Sprintf("failed to call a handler: %s", e.Reason)
}
