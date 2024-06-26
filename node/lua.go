package node

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

func Constructor(env vela.Environment) {
	xEnv = env

	ssc = newNode()
	ssc.define(env)
	env.Set("node", lua.NewFunction(newLuaNode))
}
