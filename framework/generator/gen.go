package generator

import (
	"errors"
	"fmt"
	codegen2 "github.com/RomanIschenko/notify/framework/codegen"
	gomod2 "github.com/RomanIschenko/notify/framework/gomod"
	hmapper2 "github.com/RomanIschenko/notify/framework/hmapper"
	ioc2 "github.com/RomanIschenko/notify/framework/ioc"
	parser2 "github.com/RomanIschenko/notify/framework/parser"
	"github.com/RomanIschenko/notify/framework/runtime"
	"go/format"
	"log"
)

func initDI(file *codegen2.File, imports *codegen2.ImportsDeclaration, containerName string, deps []parser2.Dependency) codegen2.Lines {
	entityDecl := codegen2.NewStruct("consumer")
	entityDecl.Add("value", "interface{}")
	entityDecl.Add("mapping", "map[string]string")
	file.Add(entityDecl)
	lines := codegen2.Lines{}

	lines.Add("consumers := []consumer{")

	for _, t := range deps {

		p, err := gomod2.RelFilePathToModPath(t.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		label := imports.GetOrCreate(p)

		mapping := `map[string]string{`
		for field, lab := range t.Dependencies {
			mapping += fmt.Sprintf(`"%s": "%s",`, field, lab.Label)
		}
		mapping += "}"
		lines.Add(fmt.Sprintf("{value: &%s.%s{}, mapping: %s},", label, t.Name, mapping))
	}
	lines.Add("}")
	lines.Add(fmt.Sprintf("for _, c := range consumers { %s.InitConsumer(c.value, c.mapping) }", containerName))
	return lines
}

func initHandlers(file *codegen2.File, imports *codegen2.ImportsDeclaration, hmapper string, ioc string, handlers []parser2.Handler) codegen2.Lines {
	lines := codegen2.Lines{}

	handlerDecl := codegen2.NewStruct("handler")
	handlerDecl.Add("example", "interface{}")
	handlerDecl.Add("prefix", "string")
	handlerDecl.Add("eventsMapping", "map[string]string")
	file.Add(handlerDecl)

	reflectModName := imports.GetOrCreate("reflect")

	lines.Add("handlers := []handler{")
	for _, handler := range handlers {
		p, err := gomod2.RelFilePathToModPath(handler.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		p = imports.GetOrCreate(p)
		mapping := "map[string]string{"
		for ev, method := range handler.EventListeners {
			mapping += fmt.Sprintf("\"%s\": \"%s\",", ev, method)
		}
		mapping += "}"
		lines.Add(fmt.Sprintf("{example: &%s.%s{}, eventsMapping: %s, prefix: \"%s\"},", p, handler.TypeName(), mapping, handler.Prefix))
	}
	lines.Add("}")

	lines.Add("for _, h := range handlers {")
	lines.Add(fmt.Sprintf("existingH, ok := %s.FindConsumer(%s.TypeOf(h.example))", ioc, reflectModName))
	lines.Add(fmt.Sprintf("var handlerValue interface{}"))
	lines.Add(fmt.Sprintf("if ok {handlerValue = existingH} else {handlerValue = h.example}"))
	lines.Add(fmt.Sprintf("%s.AddHandler(handlerValue, h.prefix, h.eventsMapping)", hmapper))
	lines.Add("}")
	return lines
}

func initConfigurators(imports *codegen2.ImportsDeclaration, builder string, cfgs []parser2.Configurator) codegen2.Lines {
	var lines codegen2.Lines

	arrDecl := "configurators := []interface{}{"

	for _, cfg := range cfgs {
		p, err := gomod2.RelFilePathToModPath(cfg.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		p = imports.GetOrCreate(p)

		arrDecl += fmt.Sprintf("%s.%s, ", p, cfg.Name)
	}

	arrDecl += "}"

	lines.Add(arrDecl)

	genRuntimeModName := imports.GetOrCreate(runtime.ImportPath)

	lines.Add(fmt.Sprintf("%s.Configure(%s, configurators...)", genRuntimeModName, builder))

	return lines
}

func initStarters(imports *codegen2.ImportsDeclaration, app string, starters []parser2.Starter) codegen2.Lines {
	var lines codegen2.Lines

	if len(starters) == 0 {
		return lines
	}

	s := starters[0]

	p, err := gomod2.RelFilePathToModPath(s.Path)
	if err != nil {
		log.Println(err)
		return lines
	}
	p = imports.GetOrCreate(p)

	genRuntimeModName := imports.GetOrCreate(runtime.ImportPath)

	lines.Add(fmt.Sprintf("%s.Start(%s, %s.%s)", genRuntimeModName, app, p, s.Name))

	return lines
}

func initHooks(file *codegen2.File, imports *codegen2.ImportsDeclaration, app string, hooks []parser2.Hook) codegen2.Lines {
	var lines codegen2.Lines

	rtModName := imports.GetOrCreate(runtime.ImportPath)



	lines.Add(fmt.Sprintf("hookContainer := %s.NewHookContainer(%s)", rtModName, app))

	hooksArrDecl := fmt.Sprintf("hookContainer.Init([]%s.Hook{", rtModName)

	for _, hook := range hooks {
		p, err := gomod2.RelFilePathToModPath(hook.Path)
		if err != nil {
			log.Println(err)
			return lines
		}
		p = imports.GetOrCreate(p)
		hooksArrDecl += fmt.Sprintf("{Event: \"%s\", Handler: %s.%s},", hook.Name, p, hook.FuncName)
	}
	hooksArrDecl += "})"

	lines.Add(hooksArrDecl)
	lines.Add("defer hookContainer.Close()")
	return lines
}

func Generate(info parser2.Info) (string, error) {

	if len(info.Starters) > 1 {
		return "", errors.New("you can't have more than one starter")
	}
	file := codegen2.File{
		Package: "main",
	}

	imports := codegen2.NewImports()
	file.Add(imports)

	mainFunc := codegen2.NewFunc("main")
	file.Add(mainFunc)

	iocModName := imports.GetOrCreate(ioc2.ImportPath)
	hmapperModName := imports.GetOrCreate(hmapper2.ImportPath)
	notifyModName := imports.GetOrCreate("github.com/RomanIschenko/notify")
	builderModName := imports.GetOrCreate("github.com/RomanIschenko/notify/framework/builder")

	mainFunc.AddLine(fmt.Sprintf("iocContainer := %s.New()", iocModName))
	mainFunc.AddLine(fmt.Sprintf("appBuilder := %s.New(iocContainer)", builderModName))
	mainFunc.AddLines(initConfigurators(imports, "appBuilder", info.Configurators))
	mainFunc.AddLine(fmt.Sprintf("app := %s.New(*appBuilder.AppConfig())", notifyModName))
	mainFunc.AddLine(fmt.Sprintf("iocContainer.Add(%s.NewEntry(app, \"\"))", iocModName))
	mainFunc.AddLine(fmt.Sprintf("handlersMapper := %s.New(app)", hmapperModName))
	mainFunc.AddLines(initHooks(&file, imports, "app", info.Hooks))
	mainFunc.AddLines(initDI(&file, imports, "iocContainer", info.Dependencies))
	mainFunc.AddLines(initHandlers(&file, imports, "handlersMapper", "iocContainer", info.Handlers))
	mainFunc.AddLines(initStarters(imports, "app", info.Starters))

	bts, err := format.Source([]byte(file.String()))

	return string(bts), err
}
