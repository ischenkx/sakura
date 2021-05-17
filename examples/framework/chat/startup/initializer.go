package startup

import (
	"fmt"
	info2 "github.com/ischenkx/notify/framework/info"
)

// notify:initialize
func Initializer(info info2.Info) {
	fmt.Println("initialized:", info.App().ID())
}
