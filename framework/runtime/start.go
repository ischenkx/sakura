package runtime

import (
	"github.com/ischenkx/notify"
	"github.com/ischenkx/notify/internal/utils"
	"reflect"
)

func Start(app *notify.App, runner interface{}) {
	val := reflect.ValueOf(runner)
	if !val.IsValid() {
		panic("runtime error: invalid starter")
	}

	typ := val.Type()

	if typ.Kind() != reflect.Func {
		panic("runtime error: starter is not a function")
	}

	in := typ.NumIn()

	if in != 1 {
		panic("runtime error: starter must accept one parameter - app")
	}

	if !utils.CompareTypes(typ.In(0), reflect.TypeOf(&notify.App{})) {
		panic("runtime error: starter must accept one parameter - app")
	}
	val.Call([]reflect.Value{reflect.ValueOf(app)})
}
