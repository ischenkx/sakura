package channels

import "strings"

func FromUser(user string) string {
	return "user/" + user
}

func FromTopic(topic string) string {
	return "topic/" + topic
}

func ParseUser(channel string) string {
	return strings.TrimPrefix(channel, "user/")
}

func ParseTopic(channel string) string {
	return strings.TrimPrefix(channel, "topic/")
}
