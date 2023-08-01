package logger

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/vela"
)

func (z *zapState) startup(ctx *fasthttp.RequestCtx) error {
	out, w := auxlib.Stdout()
	defer w.Close()

	out("start logger config change ...")
	body := ctx.Request.Body()
	var cfg config
	err := json.Unmarshal(body, &cfg)
	if err != nil {
		out("start logger config json decode fail %s", string(body))
		return err
	}

	out("reload logger config %#v", cfg)
	z.reload(newZapState(&cfg))
	out("start logger config reload succeed")
	return nil
}

func (z *zapState) define(env vela.Environment) {
	r := env.R()
	r.POST("/api/v1/inline/agent/logger", env.Then(z.startup))
}
