package parser

import (
	"fmt"
	"github.com/RomanIschenko/notify/framework/command"
	"github.com/RomanIschenko/notify/framework/common"
	"go/doc"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

type Handler struct {
	typ    *doc.Type
	Path   string
	Prefix string
	// event => method
	EventListeners map[string]string
}

func (h *Handler) TypeName() string {
	return h.typ.Name
}

func (h *Handler) Log(prefix string) {
	fmt.Printf("%sPrefix: %s\n", prefix, h.Prefix)
	for event, method := range h.EventListeners {
		fmt.Printf("%s%s => %s\n", prefix, event, method)
	}
}

func (h *Handler) Init(fileset *token.FileSet, msgstack *common.MessageStack) {
	for _, method := range h.typ.Methods {
		command.IterText(method.Doc, func(cmd command.Command, err error) {
			if err != nil {
				file := fileset.File(method.Decl.Pos())
				line := file.Line(method.Decl.Pos())
				msgstack.Error(fmt.Sprintf("failed to parse command: %s\n\tat %s:%d", err, file.Name(), line))
				return
			}

			switch cmd.Command {
			case "on":
				if ev, ok := cmd.FindFlag("event"); ok {
					h.EventListeners[ev] = method.Name
				}
			}
		})
	}
}

func getPrefixFromPath(projectFolder, path string) string {
	absPF, err := filepath.Abs(projectFolder)

	if err != nil {
		log.Println("failed to get abs path - prefix is weird now 1")
		return ""
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		log.Println("failed to get abs path - prefix is weird now 2")
		return ""
	}

	realPath := strings.TrimLeft(strings.TrimPrefix(absPath, absPF), "\\")

	return strings.Join(strings.Split(realPath, "\\"), ".")
}

func newHandler(projectFolder, path string, typ *doc.Type, flags []command.Flag) Handler {
	h := Handler{
		typ:            typ,
		Path:           path,
		Prefix:         getPrefixFromPath(projectFolder, path),
		EventListeners: map[string]string{},
	}
	for _, flag := range flags {
		switch flag.Name {
		case "prefix":
			h.Prefix = flag.Value
		}
	}
	return h
}