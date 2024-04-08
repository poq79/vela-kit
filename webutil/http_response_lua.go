package webutil

import (
	strutil "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

type ResponseL struct {
	resp *http.Response
	err  error
}

func (r *ResponseL) Text() string {
	if r.err != nil {
		return r.err.Error()
	}

	if r.resp == nil {
		return ""
	}

	chunk, err := httputil.DumpResponse(r.resp, true)
	if err != nil {
		return r.err.Error()
	}

	return strutil.B2S(chunk)
}

func (r *ResponseL) ReadBody() []byte {
	if r.err != nil || r.resp == nil {
		return nil
	}
	defer r.resp.Body.Close()

	body, err := io.ReadAll(r.resp.Body)
	if err != nil {
		return nil
	}

	return body
}

func (r *ResponseL) H(key string) string {
	if r.resp == nil {
		return ""
	}

	return r.resp.Header.Get(key)
}

func (r *ResponseL) String() string                         { return r.Text() }
func (r *ResponseL) Type() lua.LValueType                   { return lua.LTObject }
func (r *ResponseL) AssertFloat64() (float64, bool)         { return 0, false }
func (r *ResponseL) AssertString() (string, bool)           { return "", false }
func (r *ResponseL) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (r *ResponseL) Peek() lua.LValue                       { return r }

func (r *ResponseL) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "text":
		return lua.S2L(r.Text())
	case "status":
		if r.resp != nil {
			return lua.LInt(r.resp.StatusCode)
		}
		return lua.LInt(0)
	case "body":
		return lua.B2L(r.ReadBody())
	}

	if strings.HasPrefix(key, "http_") {
		return lua.S2L(r.H(strings.TrimPrefix(key, "http_")))
	}

	return lua.LNil
}

func NewResponseL(resp *http.Response, err error) *ResponseL {
	return &ResponseL{resp: resp, err: err}
}
