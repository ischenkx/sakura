package subscription

import (
	"sakura/common/data"
)

type Storage data.Storage[Subscription, Selector]

type Selector struct {
	User  *string
	Topic *string
}
