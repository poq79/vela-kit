package rtable

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/webutil"
	"net/http"
	"path/filepath"
	"reflect"
)

var ref uint32 = 0

var typeof = reflect.TypeOf((*LHandle)(nil)).String()

type config struct {
	uri     string
	method  string
	recycle bool
}

type LHandle struct {
	hctx *webutil.HttpContext
	lua.SuperVelaData
	cfg      config
	Callback func(*fasthttp.RequestCtx)
	trr      *TnlRouter
	co       *lua.LState
}

func (lh *LHandle) Name() string {
	return fmt.Sprintf("%s %s", lh.cfg.method, lh.cfg.uri)
}

func (lh *LHandle) Type() string {
	return typeof
}

func (lh *LHandle) Start() error {
	if lh.Callback == nil {
		return fmt.Errorf("router start fail not found handle")
	}
	return lh.trr.Handle(lh.cfg.method, lh.cfg.uri, lh.Callback)
}

func (lh *LHandle) Close() error {
	if lh.cfg.recycle {
		return nil
	}

	lh.cfg.recycle = true
	lh.trr.Undo(lh.cfg.method, lh.cfg.uri)
	return nil
}

func (lh *LHandle) ToCall(L *lua.LState) int {

	uri := L.IsString(1)
	handle := L.Get(2)
	if len(uri) < 2 {
		L.RaiseError("invalid router uri got empty")
		return 0
	}

	if L.CodeVM() == "" {
		L.RaiseError("not allow add router by not task")
		return 0
	}

	path := filepath.Join("/api/v1/arr/lua/", L.CodeVM(), uri)

	lh.cfg.uri = filepath.ToSlash(path)

	vda := L.NewVelaData(lh.Name(), typeof)
	if vda.IsNil() {
		vda.Set(lh)
	} else {
		old := vda.Data.(*LHandle)
		old.Close()
		vda.Set(lh)
	}

	switch handle.Type() {
	case lua.LTNil:
		lh.Callback = func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("not found handle")
		}
	case lua.LTString:
		lh.Callback = func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString(handle.String())
		}

	case lua.LTFunction:
		fn, ok := handle.AssertFunction()
		if !ok {
			L.RaiseError("invalid function value")
			return 0
		}

		pn := xEnv.P(fn)
		lh.Callback = func(ctx *fasthttp.RequestCtx) {
			co := xEnv.Clone(lh.co)
			defer xEnv.Free(co)
			co.SetValue(webutil.WEB_CONTEXT_KEY, ctx)

			err := co.CallByParam(pn, lh.hctx)
			if err != nil {
				ctx.Error(err.Error(), http.StatusInternalServerError)
				return
			}
		}

	default:
		lh.Callback = func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString(handle.String())
		}
	}

	err := lh.Start()
	if err != nil {
		L.RaiseError("inject router error %v", err)
	}

	return 0
}

func (lh *LHandle) LFunc(L *lua.LState) *lua.LFunction {
	return L.NewFunction(lh.ToCall)
}

func (trr *TnlRouter) newHandleL(L *lua.LState, method string) lua.LValue {
	lh := &LHandle{
		hctx: webutil.NewContext(),
		cfg: config{
			method: method,
		},
		trr: trr,
		co:  xEnv.Clone(L),
	}

	return lh.LFunc(L)
}
