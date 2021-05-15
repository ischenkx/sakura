package history

type PointerStorage struct {
	data map[*History]SnapshotInfo
}

func (s PointerStorage) Push(ptr Pointer) {
	if ptr.h == nil {
		return
	}
	if _, ok := s.data[ptr.h]; ok {
		return
	}
	s.data[ptr.h] = ptr.info
}

// if f returns true the history is deleted
func (s PointerStorage) Load(f func(h *History, info SnapshotInfo) bool) {
	if f == nil {
		return
	}
	for h, info := range s.data {
		if f(h, info) {
			delete(s.data, h)
		}
	}
}

func NewPointerStorage() PointerStorage {
	return PointerStorage{data: map[*History]SnapshotInfo{}}
}
