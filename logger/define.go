package logger

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/fileutil"
	"github.com/vela-ssoc/vela-kit/stdutil"
	"github.com/vela-ssoc/vela-kit/vela"
	"path/filepath"
	"sort"
)

func (z *zapState) startup(ctx *fasthttp.RequestCtx) error {
	out := stdutil.New(stdutil.Console())
	defer func() {
		_ = out.Close()
	}()

	out.Info("start logger config change ...")
	body := ctx.Request.Body()
	var cfg config
	err := json.Unmarshal(body, &cfg)
	if err != nil {
		out.ERR("start logger config json decode fail %s", string(body))
		return err
	}

	out.Info("reload logger config %#v", cfg)
	z.reload(newZapState(&cfg))
	out.Info("start logger config reload succeed")
	return nil
}

func (z *zapState) CleanAPI(ctx *fasthttp.RequestCtx) error {
	_ = z.rotate()

	dir := filepath.Dir(z.cfg.Filename)
	ignore := func(attr fileutil.Attr) bool {
		if attr.Dir || attr.Filename == z.cfg.Filename {
			return true
		}
		return false
	}

	errFn := func(err error) {
		xEnv.Error(err)
	}

	attrs := fileutil.Glob(filepath.Join(dir, "*.log"), ignore, errFn)
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].MTime.Unix() < attrs[i].MTime.Unix()
	})

	//
	//保留两个文件  但是如果单个文大于200M
	space := ctx.QueryArgs().GetUintOrZero("space")

	return fileutil.Clean(attrs, 2, 200*1024*1024, int64(space))
}

func (z *zapState) define(env vela.Environment) {
	r := env.R()
	r.POST("/api/v1/inline/agent/logger", env.Then(z.startup))
	r.GET(("/clean/agent/logger"), env.Then(z.CleanAPI))
}
