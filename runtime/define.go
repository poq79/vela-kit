package runtime

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/vela-ssoc/vela-kit/vela"
	"net/http/pprof"
)

func Select(key string, size int, reverse bool) ([]byte, error) {
	db := xEnv.Storm("runtime", key)

	var records []Record
	var err error
	if reverse {
		err = db.Select().Limit(size).Reverse().OrderBy("ID").Find(&records)
	} else {
		err = db.Select().Limit(size).OrderBy("ID").Find(&records)
	}

	//err = db.Select().Limit(size).Find(&records)
	if err != nil {
		return nil, err
	}

	chunk, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}

	return chunk, nil

}

/*
func helper(key string, ctx *fasthttp.RequestCtx) error {
	chunk, err := Select(key, Size(ctx, 100), true)
	if err != nil {
		return err
	}
	ctx.Write(chunk)
	return nil
}
*/

func helper(r *ring.Ring, ctx *fasthttp.RequestCtx) error {
	var result []interface{}
	r.Do(func(v interface{}) {
		if v == nil {
			return
		}
		result = append(result, v)
	})

	chunk, err := json.Marshal(result)
	if err != nil {
		return err
	}

	ctx.Write(chunk)
	return nil
}

func (m *monitor) define(r vela.Router) {
	r.GET("/api/v1/arr/agent/runtime/os_cpu", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		//return helper("os.cpu", ctx)
		return helper(rb.os.cpu, ctx)
	}))

	r.GET("/api/v1/arr/agent/runtime/os_mem", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		//return helper("os.mem", ctx)
		return helper(rb.os.mem, ctx)
	}))

	r.GET("/api/v1/arr/agent/runtime/agent_cpu", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		//return helper("agent.cpu", ctx)
		return helper(rb.agent.cpu, ctx)
	}))

	r.GET("/api/v1/arr/agent/runtime/agent_mem", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		return helper(rb.agent.mem, ctx)
	}))

	r.GET("/api/v1/arr/pprof/index", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index)(ctx)
		return nil
	}))

	r.GET("/api/v1/arr/pprof/cmdline", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline)(ctx)
		return nil
	}))

	r.GET("/api/v1/arr/pprof/profile", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile)(ctx)
		return nil
	}))

	r.GET("/api/v1/arr/pprof/symbol", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol)(ctx)
		return nil
	}))

	r.GET("/api/v1/arr/pprof/trace", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace)(ctx)
		return nil
	}))
	r.GET("/api/v1/arr/pprof/{name:*}", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		uv := ctx.UserValue("name")
		name, ok := uv.(string)
		if !ok {
			return fmt.Errorf("not found name")
		}
		fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler(name).ServeHTTP)(ctx)
		return nil
	}))
}
