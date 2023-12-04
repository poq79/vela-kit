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

func (i *Info) catL(L *lua.LState) int {
	if i.Entry.IsZip() {
		L.Push(lua.LNil)
		L.Pushf("%s is zip", i.Entry.File())
		return 2
	}

	if i.Entry.IsNull() {
		L.Push(lua.LNil)
		L.Pushf("%s is empty", i.Entry.File())
		return 2
	}

	chunk, err := i.Cat()
	if err != nil {
		L.Push(lua.LNil)
		L.Pushf("%s read fail %v", i.Entry.File(), err)
		return 2
	}

	L.Push(lua.B2L(chunk))
	return 1
}

func (i *Info) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ok":
		return lua.LBool(i.Err == nil)
	case "json":
		return lua.NewFunction(i.BindJSON)
	case "fastjson":
		return lua.NewFunction(i.BindFastJSON)
	case "xml":
		return lua.NewFunction(i.BindXML)
	case "line":
		return lua.NewFunction(i.BindLine)
	case "cat":
		return lua.NewFunction(i.catL)
	case "file":
		return lua.S2L(i.Entry.File())
	default:
		if i.Err == nil {
			return i.Entry.Index(L, key)
		}
		return lua.LNil
	}
}
