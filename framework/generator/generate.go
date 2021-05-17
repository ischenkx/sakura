package generator

import (
	"errors"
	"fmt"
	codegen "github.com/ischenkx/notify/framework/codegen"
	gomod2 "github.com/ischenkx/notify/framework/gomod"
	hmapper2 "github.com/ischenkx/notify/framework/hmapper"
	info2 "github.com/ischenkx/notify/framework/info"
	ioc "github.com/ischenkx/notify/framework/ioc"
	parser2 "github.com/ischenkx/notify/framework/parser"
	"github.com/ischenkx/notify/framework/runtime"
	"go/format"
	"log"
)

func initDI(file *codegen.File, imports *codegen.ImportsDeclaration, containerName string, deps []parser2.Dependency) codegen.Lines {

	iocModName := imports.GetOrCreate(ioc.ImportPath)

	var lines codegen.Lines

	arrDecl := fmt.Sprintf("var consumers = []%s.Consumer{\n", iocModName)

	for _, t := range deps {
		p, err := gomod2.RelFilePathToModPath(t.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		label := imports.GetOrCreate(p)

		mapping := `map[string]string{`
		for field, lab := range t.Dependencies {
			mapping += fmt.Sprintf("\"%s\": \"%s\",", field, lab.Label)
		}
		mapping += "}"
		arrDecl += fmt.Sprintf("{Object: &%s.%s{}, DependencyMapping: %s},\n", label, t.Name, mapping)
	}
	arrDecl += "}"

	file.Add((&codegen.Lines{}).Add(arrDecl))

	lines.Add(fmt.Sprintf("for _, c := range consumers { %s.AddConsumer(c) }", containerName))
	return lines
}

func initHandlers(file *codegen.File, imports *codegen.ImportsDeclaration, hmapper string, ioc string, handlers []parser2.Handler) codegen.Lines {
	lines := codegen.Lines{}

	handlerDecl := codegen.NewStruct("handler")
	handlerDecl.Add("example", "interface{}")
	handlerDecl.Add("prefix", "string")
	handlerDecl.Add("eventsMapping", "map[string]string")
	file.Add(handlerDecl)

	reflectModName := imports.GetOrCreate("reflect")

	arrDecl := "var handlers = []handler{\n"
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
		arrDecl += fmt.Sprintf("{example: &%s.%s{}, eventsMapping: %s, prefix: \"%s\"},\n", p, handler.TypeName(), mapping, handler.Prefix)
	}
	arrDecl += "}"

	file.Add((&codegen.Lines{}).Add(arrDecl))

	lines.Add("for _, h := range handlers {")
	lines.Add(fmt.Sprintf("handlerValue := h.example"))
	lines.Add(fmt.Sprintf("if existingHandler, ok := %s.FindConsumer(%s.TypeOf(h.example)); ok {handlerValue = existingHandler}", ioc, reflectModName))
	lines.Add(fmt.Sprintf("%s.AddHandler(handlerValue, h.prefix, h.eventsMapping)", hmapper))
	lines.Add("}")
	return lines
}

//func isPrivate(name string) bool {
//	return string(name[0]) == strings.ToLower(string(name[0]))
//}
//
//func generateNonPrivateFunc() {
//
//}

func initConfigurators(file *codegen.File, imports *codegen.ImportsDeclaration, builder string, cfgs []parser2.Configurator) codegen.Lines {
	var lines codegen.Lines

	arrDecl := "var configurators = []interface{}{\n"

	for _, cfg := range cfgs {
		p, err := gomod2.RelFilePathToModPath(cfg.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		p = imports.GetOrCreate(p)
		arrDecl += fmt.Sprintf("%s.%s,\n", p, cfg.Name)
	}

	arrDecl += "}"

	file.Add((&codegen.Lines{}).Add(arrDecl))

	genRuntimeModName := imports.GetOrCreate(runtime.ImportPath)

	lines.Add(fmt.Sprintf("%s.Configure(%s, configurators...)", genRuntimeModName, builder))

	return lines
}

func initStarters(imports *codegen.ImportsDeclaration, app string, starters []parser2.Starter) codegen.Lines {
	var lines codegen.Lines

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

func initHooks(file *codegen.File, imports *codegen.ImportsDeclaration, app string, hooks []parser2.Hook) codegen.Lines {
	var lines codegen.Lines

	if len(hooks) == 0 {
		return lines
	}

	rtModName := imports.GetOrCreate(runtime.ImportPath)

	hooksArrDecl := fmt.Sprintf("var hookList = []%s.Hook{\n", rtModName)
	for _, hook := range hooks {
		p, err := gomod2.RelFilePathToModPath(hook.Path)
		if err != nil {
			log.Println(err)
			return lines
		}
		p = imports.GetOrCreate(p)
		hooksArrDecl += fmt.Sprintf("{Event: \"%s\", Handler: %s.%s, Priority: %d},\n", hook.Name, p, hook.FuncName, hook.Priority)
	}
	hooksArrDecl += "}"

	file.Add((&codegen.Lines{}).Add(hooksArrDecl))
	lines.Add(fmt.Sprintf("hookContainer := %s.NewHookContainer(%s)", rtModName, app))
	lines.Add("hookContainer.Init(hookList)")
	lines.Add("defer hookContainer.Close()")
	return lines
}

func initInitializers(file *codegen.File, imports *codegen.ImportsDeclaration, info string, ins []parser2.Initializer) codegen.Lines {
	var lines codegen.Lines

	if len(ins) == 0 {
		return lines
	}

	arrDecl := "var initializers = []interface{}{\n"

	for _, in := range ins {
		p, err := gomod2.RelFilePathToModPath(in.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		p = imports.GetOrCreate(p)

		arrDecl += fmt.Sprintf("%s.%s,\n", p, in.Name)
	}

	arrDecl += "}"

	file.Add((&codegen.Lines{}).Add(arrDecl))

	runtimeModName := imports.GetOrCreate(runtime.ImportPath)

	lines.Add(fmt.Sprintf("%s.Initialize(%s, initializers...)", runtimeModName, info))

	return lines
}

func Generate(info parser2.Info) (string, error) {

	if len(info.Starters) > 1 {
		return "", errors.New("you can't have more than one starter")
	}
	file := codegen.File{
		Package: "main",
	}

	imports := codegen.NewImports()
	file.Add(imports)

	mainFunc := codegen.NewFunc("main")

	iocModName := imports.GetOrCreate(ioc.ImportPath)
	hmapperModName := imports.GetOrCreate(hmapper2.ImportPath)
	notifyModName := imports.GetOrCreate("github.com/ischenkx/notify")
	builderModName := imports.GetOrCreate("github.com/ischenkx/notify/framework/builder")
	infoModName := imports.GetOrCreate(info2.ImportPath)

	mainFunc.AddLine(fmt.Sprintf("iocContainer := %s.New()", iocModName))
	mainFunc.AddLine(fmt.Sprintf("appBuilder := %s.New(iocContainer)", builderModName))
	mainFunc.AddLines(initConfigurators(&file, imports, "appBuilder", info.Configurators))
	mainFunc.AddLine(fmt.Sprintf("app := %s.New(*appBuilder.AppConfig())", notifyModName))
	mainFunc.AddLine(fmt.Sprintf("iocContainer.AddEntry(%s.NewEntry(app, \"\"))", iocModName))
	mainFunc.AddLine(fmt.Sprintf("handlersMapper := %s.New(app)", hmapperModName))
	mainFunc.AddLines(initHooks(&file, imports, "app", info.Hooks))
	mainFunc.AddLine(fmt.Sprintf("runtimeInfo := %s.New(appBuilder.Context(), iocContainer, app)", infoModName))
	mainFunc.AddLines(initDI(&file, imports, "iocContainer", info.Dependencies))
	mainFunc.AddLines(initHandlers(&file, imports, "handlersMapper", "iocContainer", info.Handlers))
	mainFunc.AddLines(initInitializers(&file, imports, "runtimeInfo", info.Initializers))
	mainFunc.AddLines(initStarters(imports, "app", info.Starters))

	file.Add(mainFunc)


	bts, err := format.Source([]byte(file.String()))

	return string(bts), err
}
