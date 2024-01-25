package reflectx

import (
	"github.com/vela-ssoc/vela-kit/lua"
)

func defaultDataTypeToLValue(L *lua.LState, v interface{}) (lua.LValue, bool) {
	switch rv := v.(type) {
	case nil:
		return lua.LNil, true
	case lua.LValue:
		return rv, true

	case bool:
		return lua.LBool(rv), true
	case float64:
		return lua.LNumber(rv), true
	case float32:
		return lua.LNumber(rv), true

	case int8:
		return lua.LInt(rv), true
	case int16:
		return lua.LInt(rv), true
	case int32:
		return lua.LNumber(rv), true
	case int:
		return lua.LInt(rv), true
	case int64:
		return lua.LNumber(rv), true

	case uint8:
		return lua.LInt(rv), true
	case uint16:
		return lua.LInt(rv), true
	case uint32:
		return lua.LNumber(rv), true
	case uint:
		return lua.LInt(rv), true
	case uint64:
		return lua.LNumber(rv), true

	case []byte:
		return lua.B2L(rv), true

	case string:
		return lua.S2L(rv), true

	case []string:
		lv := v.([]string)
		n := len(lv)

		tab := L.CreateTable(n, 0)
		for i := 0; i < n; i++ {
			tab.RawSetInt(i+1, lua.S2L(lv[i]))
		}

		return tab, true

	case []interface{}:
		n := len(rv)
		tab := L.CreateTable(n, 0)
		if n == 0 {
			return tab, true
		}

		for i := 0; i < n; i++ {
			llv, ok := defaultDataTypeToLValue(L, rv[i])
			if !ok {
				return nil, false
			}
			tab.RawSetInt(i+1, llv)
		}

		return tab, true

	case map[string]interface{}:
		n := len(rv)
		tab := lua.NewMap(n, false)
		if n == 0 {
			return tab, true
		}

		for key, val := range rv {
			tab.Set(key, ToLValue(val, L))
		}

		return tab, true
	case func():
		return lua.NewFunction(func(co *lua.LState) int {
			rv()
			return 0
		}), true

	case func() error:
		return lua.NewFunction(func(co *lua.LState) int {
			if e := rv(); e != nil {
				L.Push(lua.S2L(e.Error()))
				return 1
			}
			return 0
		}), true

	case lua.LGFunction:
		return lua.NewFunction(rv), true

	case LVFace:
		return rv.ToLValue(), true

	case error:
		return lua.S2L(v.(error).Error()), true

	default:
		return nil, false
	}
}
