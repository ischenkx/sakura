package broadcaster

type messageBuffer struct {
	bts []byte
	messages []Message
}

func (mb messageBuffer) Bytes() []byte {
	return mb.bts
}

func (mb messageBuffer) Messages() []Message {
	return mb.messages
}
