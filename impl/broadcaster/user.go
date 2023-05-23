package broadcaster

import "context"

type User interface {
	ID() string
	Send(ctx context.Context, payload []byte) error
}
