package rtable

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/webutil"
)

type Client struct {
	rTab *TnlRouter
}

func (c Client) String() string                         { return "tunnel.router.client" }
func (c Client) Type() lua.LValueType                   { return lua.LTObject }
func (c Client) AssertFloat64() (float64, bool)         { return 0, false }
func (c Client) AssertString() (string, bool)           { return "", false }
func (c Client) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (c Client) Peek() lua.LValue                       { return c }

func (c Client) Exec(method string, L *lua.LState) int {
	req := L.CheckString(1)   //req url
	data := L.Get(2).String() //req body

	r, err := c.rTab.Exec(method, req, data)

	rsp := webutil.NewResponseL(r, err)
	L.Push(rsp)
	return 1
}

func (c Client) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "GET":
		return L.NewFunction(func(co *lua.LState) int {
			return c.Exec("GET", co)
		})

	case "POST":
		return L.NewFunction(func(co *lua.LState) int {
			return c.Exec("POST", co)
		})
	}

	L.RaiseError("not found %s with router client", key)
	return lua.LNil
}
