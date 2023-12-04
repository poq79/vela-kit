package webutil

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"path/filepath"
)

func HttpRedirect(L *lua.LState) int {
	n := L.GetTop()
	var path string
	var code int

	switch n {
	case 1:
		path = L.CheckString(1)
		code = 302
	case 2:
		path = L.CheckString(1)
		code = L.CheckInt(2)
	default:
		return 0
	}

	ctx := CheckMetadataCtx(L)
	ctx.Redirect(path, code)
	return 0
}

func RequestHeaderL(L *lua.LState) int {
	return HttpHeaderHelper(L, false)
}

func ResponseHeaderL(L *lua.LState) int {
	return HttpHeaderHelper(L, true)
}

func (hctx *HttpContext) Index(co *lua.LState, key string) lua.LValue {
	ctx := CheckMetadataCtx(co)
	switch key {
	case "json":
		return hctx.sayJson
	case "clone":
		return hctx.clone
	case "say":
		return hctx.say
	case "raw":
		return hctx.sayRaw
	case "file":
		return hctx.sayFile
	case "append":
		return hctx.append
	case "exit":
		return hctx.exit
	case "eof":
		return hctx.eof
	case "redirect":
		return hctx.rdt
	case "format":
		return hctx.format

	case "req_header", "rqh", "h1":
		return hctx.rqh
	case "resp_header", "rph", "h2":
		return hctx.rph
	case "try":
		return hctx.try
	case "bind":
		return hctx.bind
	case "exec":
		return hctx.exec
	}

	return K2V(ctx, key)
}

func (hctx *HttpContext) NewIndex(co *lua.LState, key string, val lua.LValue) {
	ctx := CheckMetadataCtx(co)
	switch key {
	case "path":
		ctx.URI().SetPath(val.String())
	case "root":
		ctx.SetUserValue(WEB_ROOT_HTML, filepath.Clean(val.String()))
	case "page":
		ctx.SetUserValue(WEB_DEFAULT_PAGE, filepath.Clean(val.String()))
	}

	if key == "path" {
	}
}
