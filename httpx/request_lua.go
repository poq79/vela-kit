package httpx

import (
	"bytes"
	"fmt"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"os"
	"strings"
)

func (r *Request) String() string                         { return fmt.Sprintf("web.request %p", r) }
func (r *Request) Type() lua.LValueType                   { return lua.LTObject }
func (r *Request) AssertFloat64() (float64, bool)         { return 0, false }
func (r *Request) AssertString() (string, bool)           { return "", false }
func (r *Request) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (r *Request) Peek() lua.LValue                       { return r }

func (r *Request) Index(L *lua.LState, key string) lua.LValue {
	switch key {

	case "GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH", "DELETE", "TRACE":
		r.Method = key
		return L.NewFunction(r.exec)
	case "H":
		return L.NewFunction(r.H)
	case "body":
		return L.NewFunction(r.BodyL)
	}

	return nil
}

func (r *Request) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {

	case "param", "P":
		SetQueryParam(L, r, val)

	case "header", "H":
		SetHeader(L, r, val)

	case "type", "content_type":
		r.SetHeader("content-type", val.String())

	case "body":
		r.SetBody(val.String())

	default:
		if strings.HasPrefix(key, "arg_") {
			r.SetQueryParam(key[4:], val.String())
			return
		}

		if strings.HasPrefix(key, "http_") {
			r.SetHeader(U2H(key[5:]), val.String())
			return
		}
	}
}

func (r *Request) save(L *lua.LState) int {
	n := L.GetTop()

	cover := false
	filename := L.CheckString(1)
	if n >= 2 {
		cover = L.CheckBool(2)
	}

	if _, e := os.Stat(filename); os.IsNotExist(e) || cover {
		r.SetOutput(filename)
	} else {
		r.termination = fmt.Errorf("save %s error", filename)
	}

	L.Push(r)
	return 1
}

func (r *Request) exec(L *lua.LState) int {
	if r.termination != nil {
		L.Push(NewRespE(r, r.termination))
		return 1
	}

	n := L.GetTop()
	if n <= 0 {
		return 0
	}

	uri := L.CheckString(1)
	if n == 2 {
		r.SetBody(L.Get(2).String())
	}

	if n >= 3 {
		r.SetContentType(L.CheckSockets(2))
		r.SetBody(L.Get(3).String())
	}

	res, err := r.Execute(r.Method, uri)
	if err != nil {
		L.Push(NewRespE(r, err))
		return 1
	}

	L.Push(res)
	return 1
}

func (r *Request) H(L *lua.LState) int {
	data := L.CheckString(1)
	key, val := auxlib.ParamValue(data)
	r.SetHeader(key, val)
	L.Push(r)
	return 1
}

func (r *Request) BodyL(L *lua.LState) int {
	var buf bytes.Buffer
	n := L.GetTop()
	if n == 0 {
		goto done
	}

	for i := 1; i <= n; i++ {
		buf.WriteString(L.Get(i).String())
	}

	r.SetBody(buf.Bytes())
done:
	L.Push(r)
	return 1
}
