package rtable

import (
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

type tx struct {
	ctx *fasthttp.RequestCtx
}

func (t *tx) String() string {
	return "router.ctx"
}

func (t *tx) Type() lua.LValueType {
	return lua.LTObject
}

func (t *tx) AssertFloat64() (float64, bool) {
	return 0, false
}

func (t *tx) AssertString() (string, bool) {
	return "", false
}

func (t *tx) AssertFunction() (*lua.LFunction, bool) {
	return nil, false
}

func (t *tx) Peek() lua.LValue {
	return t
}

func (t *tx) sayL(L *lua.LState) int {
	t.ctx.WriteString(L.Get(1).String())
	return 0
}

func (t *tx) exitL(L *lua.LState) int {
	t.ctx.Response.SetStatusCode(L.IsInt(1))
	return 0
}

func (t *tx) headerL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.Push(t)
		return 1
	}

	for i := 1; i <= n; i++ {
		item := L.IsString(i)
		if len(item) == 0 {
			continue
		}

		key, val := auxlib.ParamValue(item)
		t.ctx.Response.Header.Set(key, val)
	}
	L.Push(t)
	return 1
}

func (t *tx) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "say":
		return lua.NewFunction(t.sayL)
	case "header":
		return lua.NewFunction(t.headerL)
	case "exit":
		return lua.NewFunction(t.exitL)
	}

	return lua.LNil

}

func ctx2lv(ctx *fasthttp.RequestCtx) *tx {
	return &tx{ctx: ctx}
}
