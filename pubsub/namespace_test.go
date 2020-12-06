package pubsub

import (
	"fmt"
	_ "net/http/pprof"
	"testing"
)

func TestParseTopicData(t *testing.T) {
	fmt.Println(parseTopicData("chats:global"))
}

func TestTopic(t *testing.T) {

}