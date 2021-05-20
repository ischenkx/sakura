package command

type Flag struct {
	Name  string
	Value string
	emptyStr bool
	hasEq bool
}