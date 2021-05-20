package runtime

import (
	"github.com/ischenkx/notify/framework/info"
	"github.com/ischenkx/notify/internal/utils"
	"reflect"
)

func Start(info2 info.Info, runner interface{}) {
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

	if !utils.CompareTypes(typ.In(0), reflect.TypeOf(info.Info{})) {
		panic("runtime error: starter must accept one parameter - info")
	}
	val.Call([]reflect.Value{reflect.ValueOf(info2)})
}
