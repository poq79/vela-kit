package auxlib

import "github.com/vela-ssoc/vela-kit/lua"

func LToSS(L *lua.LState) []string {
	n := L.GetTop()
	if n == 0 {
		return nil
	}

	if n == 1 {
		lv := L.Get(1)
		if lv.Type() == lua.LTTable {
			return LTab2SS(lv.(*lua.LTable))
		}
	}

	var ssv []string
	for i := 1; i <= n; i++ {
		lv := L.Get(i)
		if lv.Type() == lua.LTNil {
			continue
		}
		v := lv.String()
		ssv = append(ssv, v)
	}
	return ssv
}

func S2Tab(v []string) *lua.LTable {
	tab := lua.CreateTable(len(v), 0)

	for i, item := range v {
		tab.RawSetInt(i+1, lua.S2L(item))
	}

	return tab
}
