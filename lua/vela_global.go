package lua

import (
	"fmt"
	"github.com/tidwall/pretty"
)

var _G = map[string]LValue{
	"f": NewFunction(FormatL),
	"j": NewFunction(PrettyJsonL),
}

func Format(L *LState, seek int) string {
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
			case LTString:
				args = append(args, lv.String())
			case LTBool:
				args = append(args, bool(lv.(LBool)))
			case LTNumber:
				num := lv.(LNumber)
				if num == LNumber(int(num)) {
					args = append(args, int(num))
				} else {
					args = append(args, num)
				}

			case LTInt:
				args = append(args, int(lv.(LInt)))
			case LTNil:
				args = append(args, nil)
			case LTFunction:
				args = append(args, lv)
			case LTAnyData:
				args = append(args, lv.(*AnyData).Data)
			case LTUserData:
				args = append(args, lv.(*LUserData).Value)
			case LTVelaData:
				args = append(args, lv.(*VelaData).Data)
			case LTObject:
				args = append(args, lv)
			default:
				args = append(args, lv)
			}
		}

		return fmt.Sprintf(format, args...)
	}
}

func PrettyJson(lv LValue) []byte {
	chunk := S2B(lv.String())
	return pretty.PrettyOptions(chunk, nil)
}

func FormatL(L *LState) int {
	L.Push(S2L(Format(L, 0)))
	return 1
}

func PrettyJsonL(L *LState) int {
	L.Push(B2L(PrettyJson(L.Get(1))))
	return 1
}
