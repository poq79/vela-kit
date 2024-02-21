package require

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

var (
	xEnv vela.Environment
)

func (p *pool) call(L *lua.LState) int {
	name := L.CheckString(1)
	if e := auxlib.Name(name); e != nil {
		L.RaiseError("%s invalid name", name)
		return 0
	}
	L.Push(p.RequireL(L, name))
	return 1
}

func Constructor(env vela.Environment, callback func(p Pool) error) {
	xEnv = env
	p := newPool()
	if e := callback(p); e != nil {
		xEnv.Errorf("environment constructor callback fail %v", e)
		return
	}

	_ = xEnv.Spawn(5, p.sync)

	env.Set("require", lua.NewExport("lua.require.export", lua.WithFunc(p.call)))
}
