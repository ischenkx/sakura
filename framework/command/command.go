package command


// Example
//
// notify:handler
//
// notify:inject
//
// notify:event

type Command struct {
	Command string
	Flags []Flag
}

func (cmd Command) FindFlag(name string) (string, bool) {
	for _, flag := range cmd.Flags {
		if flag.Name == name {
			return flag.Value, true
		}
	}
	return "", false
}