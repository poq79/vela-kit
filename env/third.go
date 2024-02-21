package env

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

func (env *Environment) ThirdInfo(name string) *vela.ThirdInfo {
	return env.third.Info(name)
}

func (env *Environment) Third(name string) (*vela.ThirdInfo, error) {
	return env.third.Load(name)
}

func (env *Environment) RequireL(L *lua.LState, filename string) lua.LValue {
	if env.requireHub == nil {
		return lua.LNil
	}

	return env.requireHub.RequireL(L, filename)
}

func (env *Environment) Require(filename string) lua.LValue {
	if env.requireHub == nil {
		return lua.LNil
	}

	return env.requireHub.Require(filename)
}
