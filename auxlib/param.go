package auxlib

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"strings"
)

func ParamLValue(v string) (string, lua.LValue) {
	s := strings.SplitN(v, ":", 2)
	if len(s) != 2 {
		return v, nil
	}

	key := s[0]
	val := s[1]

	bv, err := ToBoolE(val)
	if err == nil {
		return key, lua.LBool(bv)
	}

	nv, err := ToIntE(val)
	if err == nil {
		return key, lua.LInt(nv)
	}

	fv, err := ToFloat64E(val)
	if err == nil {
		return key, lua.LNumber(fv)
	}

	return s[0], lua.S2L(val)
}

func ParamValue(v string) (string, string) {
	s := strings.SplitN(v, ":", 2)
	if len(s) != 2 {
		return v, ""
	}
	return s[0], s[1]
}
