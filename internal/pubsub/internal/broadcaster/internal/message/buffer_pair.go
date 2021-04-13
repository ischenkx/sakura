package message

type BufferPair struct {
	cur     int
	buffers [2]Buffer
}

func (p *BufferPair) Buffer() *Buffer {
	return &p.buffers[p.cur]
}

func (p *BufferPair) Swap() *Buffer {
	c := p.Buffer()
	p.cur = 1 - p.cur
	return c
}

func (p *BufferPair) Close() {
	p.cur = -1
	p.buffers[0].Close()
	p.buffers[1].Close()
}

func NewBufferPair() *BufferPair {
	return &BufferPair{
		cur: 0,
		buffers: [2]Buffer{
			NewBuffer(),
			NewBuffer(),
		},
	}
}
