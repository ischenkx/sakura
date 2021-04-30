package emitter

import (
	"fmt"
	"github.com/RomanIschenko/notify"
	"github.com/RomanIschenko/notify/default/codecs"
	"github.com/RomanIschenko/notify/pubsub/message"
	"reflect"
	"sync"
)

// Emitter helps you send and receive different events
// The event encoding algorithm:
// 1) 1st byte is a length of event name
// 2) from the second to 1+<NAME_LENGTH> is the name of event
// 3) other bytes are json encoded data

