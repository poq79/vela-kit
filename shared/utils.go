package shared

import (
	"github.com/vela-ssoc/vela-kit/lua"
)

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func ToNumber(v interface{}) lua.LNumber {
	lv, ok := v.(lua.LNumber)
	if ok {
		return lv
	}

	return lua.LNumber(0)

}
