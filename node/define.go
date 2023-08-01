package node

import (
	"github.com/vela-ssoc/vela-kit/vela"
)

func (nd *node) define(env vela.Environment) {
	r := env.R()
	r.POST("/api/v1/agent/notice/upgrade", env.Then(nd.upgrade)) //升级
	r.POST("/api/v1/agent/notice/command", env.Then(nd.command)) //升级
	r.POST("/api/v1/inline/agent/node", env.Then(nd.startup))
	r.POST("/api/v1/agent/node/info", env.Then(nd.Info))
	r.POST("/api/v1/arr/agent/node/info", env.Then(nd.Info))
}
