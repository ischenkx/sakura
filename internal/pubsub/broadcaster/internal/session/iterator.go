package session

type Iterator struct {
	s   *Storage
	ids []string
	idx int
	cur *Session
}

func (iter *Iterator) Next() bool {
	if iter.idx >= len(iter.ids) {
		return false
	}
	currentId := iter.ids[iter.idx]
	iter.idx++
	if s, ok := iter.s.bucket(currentId).get(currentId); ok {
		iter.cur = s
		return true
	} else {
		return iter.Next()
	}
}

func (iter *Iterator) Get() *Session {
	return iter.cur
}
