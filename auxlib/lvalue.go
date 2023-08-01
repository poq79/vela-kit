package auxlib

import "github.com/vela-ssoc/vela-kit/lua"

func LToSS(L *lua.LState) []string {
	n := L.GetTop()
	if n == 0 {
		return nil
	}
	var ssv []string
	for i := 1; i <= n; i++ {
		lv := L.Get(i)
		if lv.Type() == lua.LTNil {
			continue
		}
		v := lv.String()
		if len(v) == 0 {
			continue
		}

		ssv = append(ssv, v)
	}
	return ssv
}
