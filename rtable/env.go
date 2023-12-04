package rtable

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

/*
vela.router.GET("/api/v1/arr/" , "hhlo")

*/

func (trr *TnlRouter) indexL(L *lua.LState, key string) lua.LValue {
	switch key {
	case "GET", "POST", "PUT", "PATCH":
		return trr.newHandleL(L, key)
	case "call":
		return lua.NewFunction(trr.callL)
	}

	return lua.LNil
}

func Constructor(env vela.Environment, callback func(rt *TnlRouter) error) {
	xEnv = env
	trr := &TnlRouter{
		hub:    make(map[string]fasthttp.RequestHandler, 32),
		router: router.New(),
	}
	trr.Listen()
	callback(trr)

	trr.GET("/api/v1/arr/agent/router/info", xEnv.Then(trr.view))
	xEnv.Set("router", lua.NewExport("lua.router.export", lua.WithIndex(trr.indexL)))
}
