package publication

type Buffer struct {
	buffer []Publication
	maxSize int
	index  int
	overrun bool
}

func (m *Buffer) Push(p Publication) {
	if m.buffer == nil {
		m.buffer = make([]Publication, 0, m.maxSize)
	}
	if len(m.buffer) >= m.maxSize {
		if m.index >= m.maxSize {
			m.index = 0
		}
		m.buffer[m.index] = p
		m.overrun = true
	} else {
		m.buffer = append(m.buffer, p)
	}
}

//Flush returns buffered publications and true if buffer was overrun
func (m *Buffer) Flush() ([]Publication, bool) {
	overrun, buffer := m.overrun, m.buffer
	m.overrun = false
	m.buffer = nil
	m.index = 0
	return buffer, overrun
}

func NewBuffer(maxSize int) Buffer {
	return Buffer{
		maxSize: maxSize,
		buffer: nil,
	}
}