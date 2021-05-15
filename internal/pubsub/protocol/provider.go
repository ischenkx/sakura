package protocol

// must be thread safe
type Provider interface {
	New() Protocol
	Put(protocol Protocol)
}
