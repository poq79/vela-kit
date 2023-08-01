package third

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

type Info struct {
	Err   error
	Entry *vela.ThirdInfo
}

func (i *Info) String() string                         { return lua.B2S(i.Byte()) }
func (i *Info) Type() lua.LValueType                   { return lua.LTObject }
func (i *Info) AssertFloat64() (float64, bool)         { return 0, false }
func (i *Info) AssertString() (string, bool)           { return i.String(), true }
func (i *Info) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (i *Info) Peek() lua.LValue                       { return i }
func (i *Info) Byte() []byte {
	if i.Err != nil {
		return nil
	}
	return i.Entry.Byte()
}

func (i *Info) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ok":
		return lua.LBool(i.Err == nil)
	default:
		if i.Err == nil {
			return i.Entry.Index(L, key)
		}

		return lua.LNil
	}
}
