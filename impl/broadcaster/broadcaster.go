package broadcaster

import (
	"context"
	"sakura"
)

type Broadcaster struct {
	sakura *sakura.Sakura
}

func (broadcaster *Broadcaster) Run(ctx context.Context) error {

}

func (broadcaster *Broadcaster) Connect(ctx context.Context, user User) error {
	
}

func (broadcaster *Broadcaster) Disconnect(ctx context.Context, id string) error {

}
