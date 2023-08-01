package logger

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

func (z *zapState) NewL(L *lua.LState) int {
	if L.CodeVM() != "startup" {
		L.RaiseError("new logger not allowed in %s", L.CodeVM())
		return 0
	}

	cfg := newConfig(L)
	//重启
	z.reload(newZapState(cfg))
	return 0
}

func (z *zapState) cleanL(L *lua.LState) int {
	z.clean()
	return 0
}

func (z *zapState) rotateL(L *lua.LState) int {
	e := z.rotate()
	if e != nil {
		L.Pushf("%v", e)
		return 1
	}

	return 0
}

func (z *zapState) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "clean":
		return lua.NewFunction(z.cleanL)
	case "rotate":
		return lua.NewFunction(z.rotateL)
	default:
		return lua.LNil
	}
}

func Constructor(env vela.Environment, callback func(vela.Log) error) {
	//初始化
	xEnv = env
	state.define(env)

	env.Set("logger", lua.NewExport("lua.logger.export", lua.WithFunc(state.NewL), lua.WithIndex(state.Index)))

	//日志配置
	if e := callback(state); e != nil {
		xEnv.Errorf("not found logger state")
	}
}
