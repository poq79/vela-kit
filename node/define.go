package node

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/fileutil"
	"github.com/vela-ssoc/vela-kit/vela"
	"path/filepath"
)

type Executable struct {
	filename string
	mtime    int64
	ctime    int64
}

func (nd *node) CleanExe(pattern string, ignore func(attr fileutil.Attr) bool, space int64) error {
	errFn := func(err error) {
		xEnv.Errorf("clean ssc executable fail %v", err)
	}
	attrs := fileutil.Glob(pattern, ignore, errFn)
	return fileutil.Clean(attrs, 2, 100*1024*1024, space)
}

func (nd *node) CleanAPI(ctx *fasthttp.RequestCtx) error {
	exe, work, backup := nd.Path()
	if len(work) == 0 {
		return fmt.Errorf("no ssc work path")
	}

	ignore := func(attr fileutil.Attr) bool {
		name := filepath.Base(attr.Filename)
		return name == "ssc-mgt.exe" || attr.Filename == exe || attr.Dir
	}

	space := ctx.QueryArgs().GetUintOrZero("space")

	errs := exception.New()
	errs.Try(work+" clean", nd.CleanExe(filepath.Join(work, "ssc-*"), ignore, int64(space)))
	errs.Try(backup+" clean", nd.CleanExe(filepath.Join(backup, "ssc-*"), ignore, int64(space)))
	return errs.Wrap()
}

func (nd *node) define(env vela.Environment) {
	r := env.R()
	_ = r.POST("/api/v1/agent/notice/upgrade", env.Then(nd.upgrade)) //升级
	_ = r.POST("/api/v1/agent/notice/command", env.Then(nd.command)) //升级
	_ = r.POST("/api/v1/agent/node/info", env.Then(nd.Info))
	_ = r.POST("/api/v1/arr/agent/node/info", env.Then(nd.Info))
	_ = r.POST("/api/v1/inline/agent/node", env.Then(nd.startup))
	_ = r.GET("/clean/agent/executable", env.Then(nd.CleanAPI))
}
