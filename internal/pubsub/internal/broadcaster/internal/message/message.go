package message

type Message struct {
	Data []byte
}

func New(data []byte) Message {
	return Message{
		Data: data,
	}
}
