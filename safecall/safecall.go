package safecall

import (
	"context"
	"time"
)

type builder struct {
	protect    bool
	timeout    time.Duration
	execFn     func() error
	completeFn func()
	errorFn    func(error)
	panicFn    func(interface{})
	timeoutFn  func()
}

func New(protect bool) *builder {
	return &builder{protect: protect}
}

func (b *builder) Timeout(timeout time.Duration) *builder {
	b.timeout = timeout
	return b
}

func (b *builder) OnComplete(fn func()) *builder {
	b.completeFn = fn
	return b
}

func (b *builder) OnError(fn func(error)) *builder {
	b.errorFn = fn
	return b
}

func (b *builder) OnPanic(fn func(interface{})) *builder {
	return b
}

func (b *builder) OnTimeout(fn func()) *builder {
	b.timeoutFn = fn
	return b
}

func (b *builder) run(handle func() error, kill context.CancelFunc, over *bool) {
	var err error
	var panicked bool

	defer func() {
		*over = true
		kill()
		if !b.protect {
			return
		}

		if cause := recover(); cause != nil {
			panicked = true
			if b.panicFn != nil {
				b.panicFn(cause)
			}
		}
	}()

	err = handle()
	if err != nil && b.errorFn != nil {
		b.errorFn(err)
		return
	}

	if err == nil && !panicked && b.completeFn != nil {
		b.completeFn()
	}
}

func (b *builder) Exec(fn func() error) {
	if fn == nil {
		return
	}

	var ctx context.Context
	var kill context.CancelFunc
	var over bool
	if b.timeout > 0 {
		ctx, kill = context.WithTimeout(context.Background(), b.timeout)
	} else {
		ctx, kill = context.WithCancel(context.Background())
	}

	over = false
	go b.run(fn, kill, &over)

	<-ctx.Done()
	if b.timeoutFn != nil && !over && ctx.Err() == context.DeadlineExceeded {
		b.timeoutFn()
	}
}
