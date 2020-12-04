package taskctx

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	go http.ListenAndServe("localhost:6565", nil)
	var arr []*TaskContext
	for i := 0; i < 100000; i++ {
		num := i + 1
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond/2*time.Duration(i))
		taskCtx := New(ctx)
		taskCtx.Defer(func() {
			fmt.Println("killing", num)
		})
		arr = append(arr, taskCtx)
	}

	for _, ctx := range arr {
		<-ctx.Done()
	}
}
