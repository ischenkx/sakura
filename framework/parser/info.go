package parser

type Info struct {
	Handlers      []Handler
	Dependencies  []Dependency
	Starters      []Starter
	Configurators []Configurator
	Hooks         []Hook
	Initializers  []Initializer
}
