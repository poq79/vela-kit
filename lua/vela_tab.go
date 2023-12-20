package lua

import "strings"

// table check int or default
func (tb *LTable) CheckInt(key string, d int) int {

	v := tb.RawGetString(key)
	switch n := v.(type) {
	case LNumber:
		return int(n)
	default:
		return d
	}
}

// table check int or default
func (tb *LTable) CheckUint32(key string, d uint32) uint32 {

	v := tb.RawGetString(key)
	switch n := v.(type) {
	case LNumber:
		return uint32(n)
	default:
		return d
	}
}

// table check string or default
func (tb *LTable) CheckString(key string, d string) string {
	v := tb.RawGetString(key)
	switch ud := v.(type) {
	case LString:
		return ud.String()
	default:
		return d
	}
}

func (tb *LTable) CheckSocket(key string, L *LState) string {
	v := tb.RawGetString(key).String()
	if e := CheckSocket(v); e != nil {
		L.RaiseError("expected socket invalid  err: %v", e)
		return ""
	}
	return v
}

func (tb *LTable) CheckSockets(key string, L *LState) string {
	v := tb.RawGetString(key).String()
	arr := strings.Split(v, ",")
	var err error

	for _, item := range arr {
		err = CheckSocket(item)
		if err != nil {
			L.RaiseError("expected %s socket error: %v ", item, err)
			return ""
		}
	}
	return v
}

func (tb *LTable) CheckVelaData(L *LState, key string) *VelaData {
	data := tb.RawGetString(key)
	if data.Type() != LTVelaData {
		L.RaiseError("invalid type , %s must be userdata , got %s", key, data.Type().String())
		return nil
	}

	return data.(*VelaData)
}

func (tb *LTable) CheckBool(key string, d bool) bool {
	v := tb.RawGetString(key)
	if lv, ok := v.(LBool); ok {
		return bool(lv)
	}
	return d
}

func (tb *LTable) IsArray() bool {
	if len(tb.array) == 0 {
		return false
	}

	if len(tb.dict) != 0 || len(tb.strdict) != 0 {
		return false
	}

	return true
}

func (tb *LTable) Range(cb func(string, LValue)) {
	if tb.strdict == nil {
		return
	}

	for k, v := range tb.strdict {
		cb(k, v)
	}
}

func (tb *LTable) Array() []LValue {
	return tb.array
}

func (tb *LTable) ToMap() *Map {
	return &Map{
		safe:  false,
		entry: tb.strdict,
	}
}

func (tb *LTable) Pickup(vt LValueType, fn func(lv LValue)) {
	arr := tb.Array()
	n := len(arr)
	for i := 0; i < n; i++ {
		item := arr[i]
		if item.Type() != vt {
			continue
		}
		fn(item)
	}
}

func (tb *LTable) Int() []int {
	var ret []int
	tb.Pickup(LTNumber, func(item LValue) {
		ret = append(ret, int(item.(LNumber)))
	})
	return ret

}

func (tb *LTable) Int64() []int64 {
	var ret []int64
	tb.Pickup(LTNumber, func(item LValue) {
		ret = append(ret, int64(item.(LNumber)))
	})
	return ret
}

func (tb *LTable) Uint() []int {
	var ret []int
	tb.Pickup(LTNumber, func(item LValue) {
		ret = append(ret, int(item.(LNumber)))
	})
	return ret

}

func (tb *LTable) Uint64() []int64 {
	var ret []int64
	tb.Pickup(LTNumber, func(item LValue) {
		ret = append(ret, int64(item.(LNumber)))
	})
	return ret
}

func (tb *LTable) Strings() []string {
	var ret []string
	tb.Pickup(LTString, func(item LValue) {
		ret = append(ret, item.String())
	})
	return ret
}

func (tb *LTable) RangeDict(cb func(LValue, LValue) bool) {
	if tb.dict == nil {
		return
	}

	for k, v := range tb.dict {
		if !cb(k, v) {
			return
		}
	}
}

func (tb *LTable) RangeStrDict(cb func(string, LValue) bool) {
	if tb.strdict == nil {
		return
	}

	for k, v := range tb.strdict {
		if !cb(k, v) {
			return
		}
	}
}
