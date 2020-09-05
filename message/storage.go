package message

type Storage interface {
	Load(appId string, ids []string) ([]Message, []string)
	Store(appId string, mes Message)
}
