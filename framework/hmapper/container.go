package hmapper

import (
	"fmt"
	"github.com/RomanIschenko/notify"
	"reflect"
)

const ImportPath = "github.com/RomanIschenko/notify/framework/hmapper"

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

			fmt.Println("Handling", absEvent)
			c.app.On(absEvent, k.Interface())
		} else {
			fmt.Println("invalid", k, method, d)
		}
	}
}


func New(app *notify.App) *Container {
	return &Container{app: app}
}