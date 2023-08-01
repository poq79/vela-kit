package third

import (
	"github.com/vela-ssoc/vela-kit/lua"
)

func (th *third) indexL(L *lua.LState, key string) lua.LValue {
	switch key {
	case "clean":
		return lua.NewFunction(th.clearThirdL)
	default:
		return lua.LNil
	}
}

func (th *third) clearThirdL(L *lua.LState) int {
	th.clear()
	return 0
}

func (th *third) loadL(L *lua.LState) int {
	name := L.CheckString(1)
	if len(name) == 0 {
		L.TypeError(1, lua.LTString)
		return 0
	}

	info, err := th.Load(name)
	if err == nil {
		L.Push(&Info{Err: nil, Entry: info})
		return 1
	}

	L.Push(&Info{Err: err, Entry: nil})
	return 1
}
