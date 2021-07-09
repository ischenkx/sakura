package swirl

type variadicIdList struct {
	count func(list IDList) int
	array func(list IDList) []string
}

func (v variadicIdList) Count() int {
	if v.count == nil {
		return len(v.Array())
	}
	return v.count(v)
}

func (v variadicIdList) Array() []string {
	if v.array == nil {
		return nil
	}
	return v.array(v)
}

type variadicSubList struct {
	count func(list SubscriptionList) int
	array func(list SubscriptionList) []Subscription
}

func (v variadicSubList) Count() int {
	if v.count == nil {
		return len(v.Array())
	}
	return v.count(v)
}

func (v variadicSubList) Array() []Subscription {
	if v.array == nil {
		return nil
	}
	return v.array(v)
}

type SubscriptionList interface {
	Count() int
	Array() []Subscription
}

type IDList interface {
	Count() int
	Array() []string
}
