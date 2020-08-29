package notify

type MessageStorage interface {
	Load(appId string, ids []string) ([]Message, []string)
	Store(appId string, mes Message)
	Clean()
}
