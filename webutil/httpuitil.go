package webutil

import (
	"bytes"
	"github.com/valyala/fasttemplate"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/strutil"
	"io"
	"os"
	"path/filepath"
)

func CheckMetadataCtx(L *lua.LState) *RequestCtx {
	v := L.Value(WEB_CONTEXT_KEY)
	if v == nil {
		L.RaiseError("invalid request context")
		return nil
	}

	ctx, ok := v.(*RequestCtx)
	if !ok {
		return nil
	}
	return ctx
}

func ReadFormL(co *lua.LState) int {
	ctx := CheckMetadataCtx(co)
	h, e := ctx.FormFile("file")
	if e != nil {
		co.RaiseError("%v", e)
		return 0
	}

	file, e := h.Open()
	if e != nil {
		co.RaiseError("%v", e)
		return 0
	}
	defer file.Close()
	name := h.Filename
	var buff bytes.Buffer
	n, e := io.Copy(&buff, file)
	if e != nil {
		co.RaiseError("%v", e)
		return 0
	}

	co.Push(lua.LString(name))
	co.Push(lua.LString(buff.String()))
	co.Push(lua.LNumber(n))
	return 3
}

//  ctx.exec({})

func ExecL(L *lua.LState) int {
	ctx := CheckMetadataCtx(L)

	html, ok := ctx.UserValue(WEB_ROOT_HTML).(string)
	if !ok {
		L.RaiseError("not found root path")
		return 0
	}

	path := L.IsString(1)
	if len(path) == 0 {
		if v, have := ctx.UserValue(WEB_DEFAULT_PAGE).(string); have {
			path = filepath.Clean(v)
		} else {
			path = "index.html"
		}
	} else {
		path = filepath.Clean(path)
	}

	file := filepath.Join(html, path)
	fd, err := os.Open(file)
	if err != nil {
		ctx.NotFound()
		return 0
	}

	text, err := io.ReadAll(fd)
	if err != nil {
		L.RaiseError("read template fail %v", err)
		return 0
	}

	entry := L.CheckIndexEx(2)

	ftt := fasttemplate.New(strutil.B2S(text), "{{", "}}")
	body := ftt.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		v := entry.Index(L, tag).String()
		return w.Write(strutil.S2B(v))
	})
	ctx.WriteString(body)
	return 0
}
