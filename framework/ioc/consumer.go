package ioc

type Consumer struct {
	// field => label
	DependencyMapping map[string]string
	Object interface{}
}
