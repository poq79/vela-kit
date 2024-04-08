package audit

import (
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/fileutil"
	"github.com/vela-ssoc/vela-kit/vela"
	"path/filepath"
	"sort"
)

func (adt *Audit) CleanAPI(ctx *fasthttp.RequestCtx) error {
	err := adt.fd.Rotate()
	if err != nil {
		return err
	}

	filename := adt.fd.Filename
	dir, name := filepath.Split(filename)
	ext := filepath.Ext(name)
	if len(ext) > 0 {
		name = name[:len(name)-len(ext)]
	}

	ignore := func(attr fileutil.Attr) bool {
		if attr.Dir { //dirctory
			return true
		}

		if attr.Filename == adt.fd.Filename {
			return true
		}

		return false
	}

	errFn := func(err error) {
		xEnv.Error(err)
	}

	pattern := filepath.Join(dir, name+"*.log")
	attrs := fileutil.Glob(pattern, ignore, errFn)

	n := len(attrs)
	if n == 0 {
		return nil
	}

	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].MTime.Unix() < attrs[j].MTime.Unix()
	})

	space := ctx.QueryArgs().GetUintOrZero("space")
	return fileutil.Clean(attrs, 2, 1024*1024*1024, int64(space)) //删除保留2 个文件， 但是最大不能超过1G
}

func (adt *Audit) define(r vela.Router) {
	_ = r.GET("/clean/audit/logger", xEnv.Then(adt.CleanAPI))
}
