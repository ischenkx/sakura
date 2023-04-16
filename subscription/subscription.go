package subscription

import "fmt"

type Subscription struct {
	User  string
	Topic string
}

func (s Subscription) ID() string {
	return fmt.Sprintf("%s_%s", s.User, s.Topic)
}
