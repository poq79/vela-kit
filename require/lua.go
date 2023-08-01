package require

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"sync"
)

var (
	once sync.Once
	xEnv vela.Environment
)

func require(L *lua.LState) int {
	name := L.CheckString(1)
	if e := auxlib.Name(name); e != nil {
		L.RaiseError("%s invalid name", name)
		return 0
	}

	L.Push(instance.require(L, name))
	return 1
}

func Constructor(env vela.Environment) {
	once.Do(func() {
		xEnv = env
		xEnv.Spawn(10, instance.sync)
		xEnv.Register(instance)
	})

	env.Set("require", lua.NewFunction(require))
}
