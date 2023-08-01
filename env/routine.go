package env

import (
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/thread"
	"sync"
	"time"
)

type routine struct {
	//coroutine lua vm 协程处理池子
	coroutine sync.Pool

	//goroutine go worker 线程池
	goroutine *thread.Pool
}

// 新建线程池子
func (env *Environment) newRoutine() {
	coroutine := sync.Pool{
		New: func() interface{} {
			return lua.NewState(lua.Options{
				CallStackSize:       lua.CallStackSize,
				RegistrySize:        lua.RegistrySize,
				IncludeGoStackTrace: env.IsDebug(),
			})
		},
	}

	goroutine, err := thread.NewPool(
		50,
		thread.WithPreAlloc(true),
	)

	exception.Fatal(err)

	env.rou = &routine{
		coroutine: coroutine,
		goroutine: goroutine,
	}
}

func (env *Environment) submit(v func()) {
	env.rou.goroutine.Submit(v)
}

func (env *Environment) Spawn(delay int, task func()) error {
	err := env.rou.goroutine.Submit(task)

	if delay != 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	return err
}

func (env *Environment) State() *lua.LState {
	return env.rou.coroutine.Get().(*lua.LState)
}

func (env *Environment) Coroutine() *lua.LState {
	return env.rou.coroutine.Get().(*lua.LState)
}

func (env *Environment) Clone(co *lua.LState) *lua.LState {
	co2 := env.rou.coroutine.Get().(*lua.LState)
	co2.Copy(co)
	return co2
}

func (env *Environment) Free(co *lua.LState) {
	if co == nil {
		return
	}

	co.Keepalive()
	env.rou.coroutine.Put(co)
}
