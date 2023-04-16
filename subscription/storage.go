package subscription

import "sakura/data"

type Storage data.Storage[Subscription, Selector]

type Selector struct {
	User  *string
	Topic *string
}
