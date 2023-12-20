package httpx

import (
	"github.com/vela-ssoc/vela-kit/lua"
)

func U2H(v string) string {
	s := lua.S2B(v)
	n := len(s)
	for i := 0; i < n; i++ {
		if s[i] == '_' {
			s[i] = '-'
		}
	}

	return lua.B2S(s)
}

func SetQueryParam(L *lua.LState, r *Request, val lua.LValue) {
	t := val.Type()
	switch t {
	case lua.LTTable:
		val.(*lua.LTable).ForEach(func(key lua.LValue, item lua.LValue) {
			r.SetQueryParam(key.String(), item.String())
		})

	default:
		r.SetQueryString(val.String())
		L.RaiseError("invalid query param type %s", t.String())
	}
}

func SetHeader(L *lua.LState, r *Request, val lua.LValue) {
	if val.Type() != lua.LTTable {
		return
	}

	tab := val.(*lua.LTable)
	tab.ForEach(func(key lua.LValue, item lua.LValue) {
		r.SetHeader(key.String(), item.String())
	})
}

func (r *Request) SetContentType(t string) {
	switch t {
	case "json", "xml", "java":
		r.SetHeader("Content-Typ", "application/"+t)
	case "gif", "jpg", "png":
		r.SetHeader("Content-Typ", "image/"+t)
	default:
		r.SetHeader("Content-Typ", t)
	}
}
