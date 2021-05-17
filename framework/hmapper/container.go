package hmapper

import (
	"fmt"
	"github.com/ischenkx/notify"
	"reflect"
)

const ImportPath = "github.com/ischenkx/notify/framework/hmapper"

type Container struct {
	app *notify.App
}

func (c *Container) AddHandler(d interface{}, prefix string, m map[string]string) {
	val := reflect.ValueOf(d)
	for ev, method := range m {
		k := val.MethodByName(method)

		if k.IsValid() {

			absEvent := ev

			if prefix != "" {
				absEvent = fmt.Sprintf("%s.%s", prefix, ev)
			}
			c.app.On(absEvent, k.Interface())
		} else {
			fmt.Println("invalid handler", k, method, d)
		}
	}
}


func New(app *notify.App) *Container {
	return &Container{app: app}
}