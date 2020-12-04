package protocol

import (
	"fmt"
)


func startsWith(s1, s2 string) bool {
	if len(s1) < len(s2) {
		return false
	}
	return s1[:len(s2)] == s2
}

type ChannelFormatter struct {
	globalPrefix, userPrefix, clientPrefix, topicPrefix string
}

func (f *ChannelFormatter) ParseChannel(ch string) (clients, users, topics []string) {
	if startsWith(ch, f.clientPrefix) {
		clients = append(clients, ch[len(f.clientPrefix)+1:])
		return
	}
	if startsWith(ch, f.userPrefix) {
		users = append(users, ch[len(f.userPrefix)+1:])
		return
	}
	if startsWith(ch, f.topicPrefix) {
		topics = append(topics, ch[len(f.topicPrefix)+1:])
		return
	}
	return
}

func (f *ChannelFormatter) AppChannel() string {
	return f.globalPrefix
}

func (f *ChannelFormatter) FmtClient(c string) string {
	return fmt.Sprintf("%s:%s", f.clientPrefix, c)
}

func (f *ChannelFormatter) FmtUser(c string) string {
	return fmt.Sprintf("%s:%s", f.userPrefix, c)
}

func (f *ChannelFormatter) FmtTopic(c string) string {
	return fmt.Sprintf("%s:%s", f.topicPrefix, c)
}

func (f *ChannelFormatter) FmtChannels(clients, users, topics []string) (chs []string) {

	length := len(users) + len(topics) + len(users)

	if length == 0 {
		return
	}

	chs = make([]string, 0, length)

	for _, c := range clients {
		chs = append(chs, f.FmtClient(c))
	}

	for _, c := range users {
		chs = append(chs, f.FmtUser(c))
	}

	for _, c := range topics {
		chs = append(chs, f.FmtTopic(c))
	}

	return chs
}

func New(appId string) *ChannelFormatter {
	return &ChannelFormatter{
		globalPrefix: fmt.Sprintf("a:%s", appId),
		userPrefix:   fmt.Sprintf("%s:u", appId),
		clientPrefix: fmt.Sprintf("%s:c", appId),
		topicPrefix:  fmt.Sprintf("%s:t", appId),
	}
}

