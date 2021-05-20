package simplehandler

import (
	"fmt"
	"github.com/ischenkx/notify/framework/helpers/hooks"
)

// notify:hook name=connect
func ConnectHook(args hooks.ConnectArgs) {
	client := args.Client
	user, _ := client.User()
	fmt.Printf("CONNECT [CLIENT=\"%s\" USER=\"%s\"]\n", args.Client.ID(), user.ID())
}

// notify:hook name=failedIncomingEvent
func FailedEventHook(args hooks.FailedIncomingEventArgs) {
	fmt.Println("failed:", args.EventError.Error)
}