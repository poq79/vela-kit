package plugin

import (
	"github.com/vela-ssoc/vela-kit/env"
	"github.com/vela-ssoc/vela-kit/lua"
	"runtime"
)

var xEnv *env.Environment

func newLuaPlugin(L *lua.LState) int {
	xEnv.Errorf("plugin not support darwin")
	return 0
}

func Constructor(env *env.Environment) {
	xEnv = env
	env.Infof("plugin running in %s", runtime.GOOS)
	env.Set("load", lua.NewFunction(newLuaPlugin))
}
