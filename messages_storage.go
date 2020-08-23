package notify

type MessagesStorage interface {
	Load(appId string, ids []string) ([]Message, []string)
	Store(appId string, mes Message)
	Clean()
}
