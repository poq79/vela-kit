package auxlib

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
)

func Format(L *lua.LState, seek int) string {
	n := L.GetTop()
	if seek > n {
		return ""
	}

	offset := n - seek
	switch offset {
	case 0:
		return ""
	case 1:
		return L.Get(seek + 1).String()
	default:
		format := L.CheckString(seek + 1)
		var args []interface{}

		for idx := seek + 2; idx <= n; idx++ {
			lv := L.Get(idx)
			switch lv.Type() {
			case lua.LTString:
				args = append(args, lv.String())
			case lua.LTBool:
				args = append(args, bool(lv.(lua.LBool)))
			case lua.LTNumber:
				num := lv.(lua.LNumber)
				if num == lua.LNumber(int(num)) {
					args = append(args, int(num))
				} else {
					args = append(args, num)
				}

			case lua.LTInt:
				args = append(args, int(lv.(lua.LInt)))
			case lua.LTNil:
				args = append(args, nil)
			case lua.LTFunction:
				args = append(args, lv)
			case lua.LTAnyData:
				args = append(args, lv.(*lua.AnyData).Data)
			case lua.LTUserData:
				args = append(args, lv.(*lua.LUserData).Value)
			case lua.LTVelaData:
				args = append(args, lv.(*lua.VelaData).Data)
			case lua.LTObject:
				args = append(args, lv)
			default:
				args = append(args, lv)
			}
		}

		return fmt.Sprintf(format, args...)
	}
}
