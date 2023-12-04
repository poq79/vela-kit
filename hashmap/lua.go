package hashmap

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/mime"
	"github.com/vela-ssoc/vela-kit/vela"
)

func newLuaHashMap(L *lua.LState) int {
	val := L.Get(1)
	var hm HMap
	switch val.Type() {
	case lua.LTTable:
		hm = New(0)
		val.(*lua.LTable).Range(func(key string, val lua.LValue) {
			hm.NewIndex(L, key, val)
		})

	default:
		n, _ := val.AssertFloat64()
		hm = New(int(n))
	}
	L.Push(hm)
	return 1
}

func Constructor(env vela.Environment) {
	xEnv = env
	mime.Register((HMap)(nil), Encode, Decode)
	xEnv.Set("hm", lua.NewExport("lua.hashmap.export", lua.WithFunc(newLuaHashMap)))
	xEnv.Set("table", lua.NewExport("lua.table.global", lua.WithFunc(newLuaHashMap)))

}
