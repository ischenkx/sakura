package codegen

import "fmt"

type File struct {
	Package string
	Parts []fmt.Stringer
}

func (f *File) Add(s fmt.Stringer) {
	f.Parts = append(f.Parts, s)
}

func (f File) String() string {
	code := fmt.Sprintf("package %s", f.Package)

	for _, part := range f.Parts {
		code += fmt.Sprintf("\n%s\n", part.String())
	}

	return code
}