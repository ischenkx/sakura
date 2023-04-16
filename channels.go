package sakura

import "fmt"

func UserChannel(id string) string {
	return fmt.Sprintf("user/%s", id)
}

func TopicChannel(id string) string {
	return fmt.Sprintf("topic/%s", id)
}
