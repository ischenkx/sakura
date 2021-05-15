package parser

type Dependency struct {
	Path string
	Name string
	Dependencies map[string]FieldInfo
}

type FieldInfo struct {
	Label string
}
