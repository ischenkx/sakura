package subscription

import "errors"

type State uint

const (
	Active State = iota + 1
	Inactive
)

var SeqInvalidErr = errors.New("sequence id is invalid")

type Subscription struct {
	state State
	seq	  int64
}

func (s *Subscription) Active() bool {
	return s.state == Active
}

func (s *Subscription) Activate(seq int64) error {
	if s.seq > seq {
		return SeqInvalidErr
	}
	s.seq = seq
	s.state = Active
	return nil
}

func (s *Subscription) Deactivate(seq int64) error {
	if s.seq > seq {
		return SeqInvalidErr
	}
	s.seq = seq
	s.state = Inactive
	return nil
}