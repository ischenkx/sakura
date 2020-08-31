package message

type Buffer struct {
	buffer []string
	maxSize int
	index  int
	overrun bool
}

func (m *Buffer) Push(id string) {
	if m.buffer == nil {
		m.buffer = []string{}
	}
	if len(m.buffer) >= m.maxSize {
		if m.index >= m.maxSize {
			m.index = 0
		}
		m.buffer[m.index] = id
		m.overrun = true
	} else {
		m.buffer = append(m.buffer, id)
	}
}

func (m *Buffer) Reset() ([]string, bool) {
	overrun, buffer := m.overrun, m.buffer
	m.overrun = false
	m.buffer = nil
	m.index = 0
	return buffer, overrun
}

func NewBuffer(maxSize int) Buffer {
	return Buffer{
		maxSize: maxSize,
	}
}