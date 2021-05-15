package runtime

import (
	"github.com/RomanIschenko/notify/framework/builder"
	"github.com/RomanIschenko/notify/internal/utils"
	"reflect"
)

const ImportPath = "github.com/RomanIschenko/notify/framework/runtime"

func Configure(b *builder.Builder, fns ...interface{}) {
	for _, f := range fns {
		val := reflect.ValueOf(f)
		if !val.IsValid() {
			panic("runtime error: configurator is invalid")
		}

		typ := val.Type()

		if typ.Kind() != reflect.Func {
			panic("runtime error: configurator is not a function")
		}

		in := typ.NumIn()

		if in != 1 {
			panic("runtime error: configurator must accept one parameter - builder")
		}

		if !utils.CompareTypes(typ.In(0), reflect.TypeOf(&builder.Builder{})) {
			panic("runtime error: configurator must accept one parameter - builder")
		}
		val.Call([]reflect.Value{reflect.ValueOf(b)})
	}
}
