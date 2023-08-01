package cidr

import "github.com/vela-ssoc/vela-kit/lua"

func Check(L *lua.LState) []*IP {
	var tt []*IP

	n := L.GetTop()
	if n == 0 {
		return tt
	}

	for i := 1; i <= n; i++ {
		tv, er := Parse(L.CheckString(i))
		if er != nil {
			L.RaiseError("invalid cidr #%d", i)
			return tt
		}
		tt = append(tt, tv)
	}

	return tt
}
