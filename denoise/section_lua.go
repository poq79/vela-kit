package denoise

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

func (sec *Section) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "ttl":
		ttl := lua.IsInt(val)
		if ttl == 0 {
			L.RaiseError("invalid ttl")
			return
		}
		sec.TTL = ttl

	case "count":
		count := lua.IsInt(val)
		if count == 0 {
			L.RaiseError("invalid count")
			return
		}
		sec.Count = count

	case "index":
		switch val.Type() {
		case lua.LTString:
			sec.Index = append(sec.Index, val.String())
		case lua.LTTable:
			sec.Index = auxlib.LTab2SS(lua.CheckTable(L, val))
		default:
			L.RaiseError("invalid index")
			return

		}
	}
}
