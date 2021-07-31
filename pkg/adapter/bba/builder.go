package bba

import (
	"github.com/ischenkx/swirl"
	"github.com/ischenkx/swirl/pkg/adapter/bba/common"
)

type Builder struct {
	inter common.PubSub
}

func (b Builder) Build(app *swirl.App) swirl.Adapter {
	return &adapter{
		app:   app,
		pubsub: b.inter,
	}
}

func NewBuilder(inter common.PubSub) swirl.AdapterBuilder {
	return Builder{inter}
}
