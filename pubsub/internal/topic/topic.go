package topic

type Topic struct {
	descriptor uint64
}

func (t *Topic) Descriptor() uint64 {
	return t.descriptor
}

func New(desc uint64) *Topic {
	return &Topic{
		descriptor:    desc,
	}
}