package codegen

import (
	"fmt"
	"path/filepath"
	"strings"
)

const nameCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func validName(n string) bool {
	for _, c := range n {
		if !strings.ContainsRune(nameCharset, c) {
			return false
		}
	}

	return true
}

type StructDeclaration struct {
	Name string
	Fields map[string]string
}

func (s *StructDeclaration) Add(n, t string) {
	s.Fields[n] = t
}

func (s *StructDeclaration) String() string {
	code := fmt.Sprintf("type %s struct {\n", s.Name)

	for name, t := range s.Fields {
		code += fmt.Sprintf("\n%s %s\n", name, t)
	}

	code += "\n}\n"

	return code
}

type ImportsDeclaration struct {
	// path => label
	labels map[string]string
}

func (d *ImportsDeclaration) labelExists(l string) bool {
	for _, l1 := range d.labels {
		if l1 == l {
			return true
		}
	}
	return false
}

func (d *ImportsDeclaration) GetOrCreate(p string) string {
	l, ok := d.labels[p]
	if !ok {
		baseName := filepath.Base(p)

		if baseName == "." || !validName(baseName) {
			baseName = RandStr(12)
		}
		counter := 0
		l = baseName

		for {
			if d.labelExists(l) {
				counter += 1
				l = fmt.Sprintf("%s%d", baseName, counter)
				continue
			}
			break
		}

		d.labels[p] = l
	}

	return l
}

func (d *ImportsDeclaration) String() string {
	if len(d.labels) == 0 {
		return ""
	}
	code := "import (\n"

	for path, label := range d.labels {
		code += fmt.Sprintf("%s \"%s\"\n", label, path)
	}

	code += "\n)\n"

	return code
}

type FuncDeclaration struct {
	Name string
	Args map[string]string
	Returns map[string]string
	Lines []string
}

func (d *FuncDeclaration) AddLines(ls Lines) {
	d.Lines = append(d.Lines, ls.lines...)
}

func (d *FuncDeclaration) AddLine(s string) {
	d.Lines = append(d.Lines, s)
}

func (d *FuncDeclaration) String() string {
	code := fmt.Sprintf("func %s(", d.Name)

	for name, typ := range d.Args {
		code += fmt.Sprintf("%s %s,", name, typ)
	}

	code += ")"

	if len(d.Returns) > 0 {
		code += " ("
		for name, typ := range d.Returns {
			code += fmt.Sprintf("%s %s,", name, typ)
		}
		code += ") "
	}

	code += "{\n"

	code += strings.Join(d.Lines, "\n")
	code += "\n}\n"

	return code
}


func NewFunc(name string) *FuncDeclaration {
	return &FuncDeclaration{
		Name:    name,
		Args: map[string]string{},
		Returns: map[string]string{},
		Lines:   nil,
	}
}

func NewStruct(name string) *StructDeclaration {
	return &StructDeclaration{
		Name:   name,
		Fields: map[string]string{},
	}
}

func NewImports() *ImportsDeclaration {
	return &ImportsDeclaration{labels: map[string]string{}}
}

type Lines struct {
	lines []string
}

func (c *Lines) Add(s string) *Lines {
	c.lines = append(c.lines, s)
	return c
}

func (c *Lines) String() string {
	return strings.Join(c.lines, "\n")
}