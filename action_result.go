package notify

import "github.com/RomanIschenko/notify/pubsub"

type ActionResult struct {
	Input interface{}
	PubsubResult pubsub.Result
}
