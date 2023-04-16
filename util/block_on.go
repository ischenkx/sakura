package util

import "context"

func BlockOn(ctx context.Context) {
	<-ctx.Done()
}
