package message

type Message struct {
	Data []byte
	Meta interface{}
}

func New(data []byte) Message {
	return Message{
		Data: data,
	}
}
