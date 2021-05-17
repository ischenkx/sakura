package runtime

import (
	info2 "github.com/ischenkx/notify/framework/info"
	"github.com/ischenkx/notify/internal/utils"
	"reflect"
)

func Initialize(info info2.Info, initializers ...interface{}) {
	for _, i := range initializers {
		val := reflect.ValueOf(i)
		if !val.IsValid() {
			panic("runtime error: initializer is invalid")
		}

		typ := val.Type()

		if typ.Kind() != reflect.Func {
			panic("runtime error: initializer is not a function")
		}

		in := typ.NumIn()

		if in != 1 {
			panic("runtime error: initializer must accept one parameter - info")
		}

		if !utils.CompareTypes(typ.In(0), reflect.TypeOf(info2.Info{})) {
			panic("runtime error: configurator must accept one parameter - info")
		}
		val.Call([]reflect.Value{reflect.ValueOf(info)})
	}
}
