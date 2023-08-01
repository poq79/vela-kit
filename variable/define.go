package variable

import (
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/vela"
)

func (hub *Hub) define(env vela.Environment) {
	r := env.R()
	r.POST("/api/v1/inline/agent/extends", func(ctx *fasthttp.RequestCtx) {

	})

}
