package pcall

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/vela"
	"gopkg.in/tomb.v2"
	"runtime"
	"time"
)

type Safe struct {
	xEnv vela.Environment
	tomb tomb.Tomb

	buffer  int
	err     error
	ok      func()
	handle  func() error
	onError func(error)
	timeout time.Duration
}

func (sa *Safe) Warp() error {
	return sa.err
}

func (sa *Safe) Buffer(n int) {
	sa.buffer = n
}

func (sa *Safe) Err(fn func(error)) *Safe {
	sa.onError = fn
	return sa
}

func (sa *Safe) Ok(fn func()) *Safe {
	sa.ok = fn
	return sa
}

func (sa *Safe) Time(tv time.Duration) *Safe {
	sa.timeout = tv
	return sa
}

func (sa *Safe) Defer() {
	if !sa.debug() {
		goto callback
	}

	defer func() {
		if cause := recover(); cause != nil {
			buf := make([]byte, sa.buffer)
			n := runtime.Stack(buf[:], false)
			sa.err = fmt.Errorf("safe func exec panic: %v  %s", cause, string(buf[:n]))
		}
	}()

callback:
	sa.tomb.Kill(sa.err)
	if sa.err == nil && sa.ok != nil {
		sa.ok()
		return
	}

	if sa.err != nil && sa.onError != nil {
		sa.onError(sa.err)
		return
	}
}

func (sa *Safe) debug() bool {
	if sa.xEnv == nil {
		return false
	}

	return sa.xEnv.Mode() == "debug"
}

func (sa *Safe) Spawn() {
	go func() {
		defer sa.Defer()

		// 执行程序
		sa.err = sa.handle()

		return
	}()

	select {

	case <-sa.tomb.Dying():
		return
	case <-time.After(sa.timeout):

		sa.xEnv.Errorf("safe exec time out %d", sa.timeout)
		return
	}
}

func (sa *Safe) Env(env vela.Environment) *Safe {
	sa.xEnv = env
	return sa
}

func Exec(v interface{}, args ...interface{}) *Safe {
	safe := &Safe{
		buffer:  8912,
		timeout: time.Millisecond,
	}

	switch fn := v.(type) {
	case func():
		safe.handle = func() error {
			fn()
			return nil
		}

	case func() error:
		safe.handle = fn

	case func(...interface{}):
		safe.handle = func() error {
			fn(args)
			return nil
		}

	case func(...interface{}) error:
		safe.handle = func() error {
			return fn(args)
		}

	default:
		safe.handle = func() error {
			return fmt.Errorf("invalid function type")
		}

	}

	return safe
}
