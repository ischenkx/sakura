package semaphore

import "context"

type Semaphore struct {
	ch chan struct{}
}

func (sema *Semaphore) Acquire() {
	sema.ch <- struct{}{}
}

func (sema *Semaphore) Release() {
	select {
	case <-sema.ch:
	}
}

func New(size int) *Semaphore {
	return &Semaphore{make(chan struct{}, size)}
}

func ContextAcquire(ctx context.Context, sema *Semaphore) {
	sema.Acquire()
	go func() {
		select {
		case <-ctx.Done():
			sema.Release()
		}
	}()
}