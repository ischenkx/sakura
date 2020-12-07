package taskctx

import (
	stdcontext "context"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type TaskContext struct {
	ctx       stdcontext.Context
	done 	  chan struct{}
	waitGroup sync.WaitGroup
	canceled bool
	deferred []func()
	mu sync.RWMutex
}

func (t *TaskContext) flushDeferred() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := len(t.deferred); i >= 0; i-- {
		f := t.deferred[i]
		f()
	}

	t.deferred = nil
}

func (t *TaskContext) newTask() bool {
	t.mu.RLock()
	isCanceled := t.canceled
	t.mu.RUnlock()
	if isCanceled {
		return false
	}
	t.waitGroup.Add(1)
	return true
}

func (t *TaskContext) completeTask() {
	t.waitGroup.Done()
}

func (t *TaskContext) Done() <-chan struct{} {
	return t.done
}

func (t *TaskContext) Err() error {
	return t.ctx.Err()
}

func (t *TaskContext) Deadline() (deadline time.Time, ok bool) {
	return t.ctx.Deadline()
}

func (t *TaskContext) Value(key interface{}) interface{} {
	return t.ctx.Value(key)
}

// Terminate kills the context
func (t *TaskContext) Terminate() {
	t.mu.Lock()
	t.canceled = true
	t.mu.Unlock()
	t.flushDeferred()
	close(t.done)
}

// like Terminate but waits for all tasks to complete
func (t *TaskContext) Cancel() {
	t.mu.Lock()
	t.canceled = true
	t.mu.Unlock()
	t.waitGroup.Wait()
	t.flushDeferred()
	close(t.done)
}

func (t *TaskContext) Do(f func()) {
	if f == nil {
		return
	}
	if t.newTask() {
		f()
		t.completeTask()
	}
}

// defer executes some action right before context is canceled
func (t *TaskContext) Defer(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.deferred = append(t.deferred, f)
}

func (t *TaskContext) Go(f func()) {
	if f == nil {
		return
	}

	t.newTask()
	go func() {
		f()
		t.completeTask()
	}()
}

func (t *TaskContext) waitCtx() {
	select {
	case <-t.done:
	case <-t.ctx.Done():
		t.Cancel()
	}
}

func New(ctx context.Context) *TaskContext {
	t := &TaskContext{
		ctx:       ctx,
		done:	make(chan struct{}),
		waitGroup: sync.WaitGroup{},
	}
	go t.waitCtx()
	return t
}