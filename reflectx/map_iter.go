package reflectx

import "github.com/vela-ssoc/vela-kit/lua"

func mapCall(L *lua.LState) int {
	ref, _ := check(L, 1)

	keys := ref.MapKeys()
	i := 0
	fn := func(L *lua.LState) int {
		if i >= len(keys) {
			return 0
		}
		L.Push(ToLValue(keys[i].Interface(), L))
		L.Push(ToLValue(ref.MapIndex(keys[i]).Interface(), L))
		i++
		return 2
	}
	L.Push(L.NewFunction(fn))
	return 1
}
