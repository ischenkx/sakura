package message

type Batch struct {
	buffer Buffer
	bytes  []byte
}

func (b Batch) Buffer() Buffer {
	return b.buffer
}

func (b Batch) Bytes() []byte {
	return b.bytes
}

func NewBatch(buf Buffer, data []byte) Batch {
	return Batch{
		buffer: buf,
		bytes:  data,
	}
}
