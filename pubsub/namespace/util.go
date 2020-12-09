package namespace

import "strings"

func parseTopicData(topic string) (id, ns string) {
	i := strings.Index(topic, ":")
	if i < 0 {
		return topic, ""
	}
	return topic, topic[:i]
}