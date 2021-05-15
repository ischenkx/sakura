package parser

import (
	"fmt"
	"github.com/RomanIschenko/notify/framework/command"
	"github.com/RomanIschenko/notify/framework/common"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

type Parser struct {
	Folder string
}

func (p *Parser) CollectInfo() (Info, common.MessageStack) {
	var (
		info     Info
		errstack common.MessageStack
	)

	err := filepath.WalkDir(p.Folder, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			errstack.Error("error while walking:"+err.Error())
			return nil
		}
		if !entry.IsDir() {
			return nil
		}
		fileSet := token.NewFileSet()
		pkgs, err := parser.ParseDir(fileSet, path, nil, parser.ParseComments)

		if err != nil {
			errstack.Error(fmt.Sprintf("error while parsing directory (%s): %s", path, err))
			return nil
		}

		for _, pkg := range pkgs {

			dc := doc.New(pkg, "", doc.AllDecls|doc.PreserveAST)

			info.Configurators = append(info.Configurators, p.collectConfigurators(path, fileSet, dc, &errstack)...)
			info.Starters = append(info.Starters, p.collectStarters(path, fileSet, dc, &errstack)...)
			info.Handlers = append(info.Handlers, p.collectHandlerInfo(path, fileSet, dc, &errstack)...)
			info.Dependencies = append(info.Dependencies, p.collectDependencies(path, fileSet, dc, &errstack)...)
			info.Hooks = append(info.Hooks, p.collectHooks(path, fileSet, dc, &errstack)...)
		}

		return nil
	})

	if err != nil {
		errstack.Error(fmt.Sprintf("error while walking directory (%s): %s", p.Folder, err))
	}

	return info, errstack
}

func (p *Parser) collectHooks(path string, fileSet *token.FileSet, dc *doc.Package, errstack *common.MessageStack) []Hook {
	var hooks []Hook

	for _, t := range dc.Funcs {
		file := fileSet.File(t.Decl.Pos())
		funcLine := file.Line(t.Decl.Pos())

		command.IterText(t.Doc, func(cmd command.Command, err error) {
			if err != nil {
				fmt.Println("failure:", t.Name)
				errstack.Error(fmt.Sprintf("error while parsing command: %s\n\tat %s:%d", err, file.Name(), funcLine))
				return
			}

			switch cmd.Command {
			case "hook":
				name, ok := cmd.FindFlag("name")

				if !ok {
					errstack.Error(
						fmt.Sprintf("failed to find a \"name\" flag for a hook\n\tat %s:%d", file.Name(), funcLine),
					)
					return
				}

				h := Hook{
					Path:     path,
					Name:     name,
					FuncName: t.Name,
				}

				hooks = append(hooks, h)
			}
		})
	}
	return hooks
}

func (p *Parser) collectConfigurators(path string, fileSet *token.FileSet, dc *doc.Package, errstack *common.MessageStack) []Configurator {
	var cfgs []Configurator

	for _, t := range dc.Funcs {
		file := fileSet.File(t.Decl.Pos())
		funcLine := file.Line(t.Decl.Pos())

		command.IterText(t.Doc, func(cmd command.Command, err error) {
			if err != nil {
				errstack.Error(fmt.Sprintf("error while parsing command: %s\n\tat %s:%d", err, file.Name(), funcLine))
				return
			}
			switch cmd.Command {
			case "config":
				cfg := Configurator{
					Path: path,
					Name: t.Name,
				}
				cfgs = append(cfgs, cfg)
			}
		})

	}

	return cfgs
}

func (p *Parser) collectStarters(path string, fileSet *token.FileSet, dc *doc.Package, errstack *common.MessageStack) []Starter {
	var starters []Starter

	for _, t := range dc.Funcs {
		file := fileSet.File(t.Decl.Pos())
		funcLine := file.Line(t.Decl.Pos())

		command.IterText(t.Doc, func(cmd command.Command, err error) {
			if err != nil {
				errstack.Error(fmt.Sprintf("error while parsing command: %s\n\tat %s:%d", err, file.Name(), funcLine))
				return
			}
			switch cmd.Command {
			case "starter":
				s := Starter{
					Path: path,
					Name: t.Name,
				}
				starters = append(starters, s)
			}
		})
	}

	return starters
}

func (p *Parser) collectDependencies(path string, fileSet *token.FileSet, dc *doc.Package, errstack *common.MessageStack) []Dependency {

	var types []Dependency

	for _, t := range dc.Types {

		file := fileSet.File(t.Decl.Pos())
		funcLine := file.Line(t.Decl.Pos())

		specs := t.Decl.Specs

		if len(specs) == 0 {
			continue
		}

		typeSpec, ok := t.Decl.Specs[0].(*ast.TypeSpec)

		if !ok {
			continue
		}

		astTyp, ok := typeSpec.Type.(*ast.StructType)

		if !ok {
			continue
		}

		fields := astTyp.Fields.List

		var typ Dependency

		typ.Path = path
		typ.Name = typeSpec.Name.Name
		typ.Dependencies = map[string]FieldInfo{}

		for _, field := range fields {
			if field.Doc == nil {
				continue
			}

			for _, com := range field.Doc.List {
				comment := strings.TrimLeft(strings.TrimLeft(com.Text, "//"), " ")
				if !command.IsExpression(comment) {
					continue
				}

				cmd, err := command.Parse(comment)

				if err != nil {
					errstack.Error(fmt.Sprintf("error while parsing command: %s\n\tat %s:%d", err, file.Name(), funcLine))

					continue
				}

				switch cmd.Command {
				case "inject":
					var label = ""
					if val, ok := cmd.FindFlag("label"); ok {
						label = val
					}
					info := FieldInfo{
						Label: label,
					}
					typ.Dependencies[field.Names[0].Name] = info
				}
			}

		}
		if len(typ.Dependencies) > 0 {
			types = append(types, typ)
		}
	}

	return types
}

func (p *Parser) collectHandlerInfo(path string, fileSet *token.FileSet, dc *doc.Package, errstack *common.MessageStack) []Handler {

	var handlers []Handler

	for _, t := range dc.Types {
		file := fileSet.File(t.Decl.Pos())
		funcLine := file.Line(t.Decl.Pos())

		command.IterText(t.Doc, func(cmd command.Command, err error) {
			if err != nil {
				errstack.Error(fmt.Sprintf("error while parsing command: %s\n\tat %s:%d", err, file.Name(), funcLine))
				return
			}
			switch cmd.Command {
			case "handler":
				handlers = append(handlers, newHandler(p.Folder, path, t, cmd.Flags))
			}
		})
	}

	return handlers
}
