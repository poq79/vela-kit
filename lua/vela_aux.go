package lua

import (
	"encoding/json"
	"errors"
	"net"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	InvalidFormat = errors.New("invalid format")
	InvalidIP     = errors.New("invalid ip addr")
	InvalidPort   = errors.New("expect check socket err: port <1 or port > 65535")
)

func ToLValue(v interface{}) LValue {
	switch converted := v.(type) {
	case nil:
		return LNil
	case bool:
		return LBool(converted)
	case float64:
		return LNumber(converted)
	case float32:
		return LNumber(converted)
	case int8:
		return LInt(converted)
	case int16:
		return LInt(converted)
	case int:
		return LInt(converted)
	case int32:
		return LInt(converted)
	case uint8:
		return LUint(converted)
	case uint16:
		return LUint(converted)
	case uint:
		return LUint(converted)
	case uint32:
		return LUint(converted)
	case int64:
		return LInt64(converted)
	case uint64:
		return LUint64(converted)
	case string:
		return S2L(converted)

	case []byte:
		return B2L(converted)

	case time.Time:
		tt := float64(converted.UTC().UnixNano()) / float64(time.Second)
		return LNumber(tt)
	case error:
		if v == nil {
			return LNil
		}
		return S2L(converted.Error())

	case []string:
		n := len(converted)

		tab := newLTable(n, 0)
		for i := 0; i < n; i++ {
			tab.RawSetInt(i+1, S2L(converted[i]))
		}
		return tab

	case []int:
		n := len(converted)

		tab := newLTable(n, 0)
		for i := 0; i < n; i++ {
			tab.RawSetInt(i+1, LInt(converted[i]))
		}
		return tab

	case func():
		return NewFunction(func(_ *LState) int {
			converted()
			return 0
		})

	case LValue:
		return converted

	case LVFace:
		return converted.ToLValue()
	}

	return NewAnyData(v)
}

func IsString(v LValue) string {
	d, ok := v.AssertString()
	if ok {
		return d
	}

	return ""
}

func IsTrue(v LValue) bool {
	if lv, ok := v.(LBool); ok {
		return bool(lv) == true
	}
	return false
}

func IsFalse(v LValue) bool {
	if lv, ok := v.(LBool); ok {
		return bool(lv) == false
	}
	return false
}

func IsNumber(v LValue) LNumber {
	if lv, ok := v.(LNumber); ok {
		return lv
	}
	return 0
}

func IsInt(v LValue) int {
	if intv, ok := v.(LNumber); ok {
		return int(intv)
	}

	if intv, ok := v.(LInt); ok {
		return int(intv)
	}

	return 0
}

func IsFunc(v LValue) *LFunction {
	fn, _ := v.(*LFunction)
	return fn
}

func IsNull(v []byte) bool {
	if len(v) == 0 {
		return true
	}
	return false
}

func CheckInt(L *LState, lv LValue) int {
	if intv, ok := lv.(LNumber); ok {
		return int(intv)
	}
	L.RaiseError("must be int , got %s", lv.Type().String())
	return 0
}

func CheckIntOrDefault(L *LState, lv LValue, d int) int {
	if intv, ok := lv.(LNumber); ok {
		return int(intv)
	}
	return d
}

func CheckInt64(L *LState, lv LValue) int64 {
	if intv, ok := lv.(LNumber); ok {
		return int64(intv)
	}
	L.RaiseError("must be int64 , got %s", lv.Type().String())
	return 0
}

func CheckNumber(L *LState, lv LValue) LNumber {
	if lv, ok := lv.(LNumber); ok {
		return lv
	}
	L.RaiseError("must be LNumber , got %s", lv.Type().String())
	return 0
}

func CheckString(L *LState, lv LValue) string {
	if lv, ok := lv.(LString); ok {
		return string(lv)
	} else if LVCanConvToString(lv) {
		return LVAsString(lv)
	}
	return ""
}

func CheckBool(L *LState, lv LValue) bool {
	if lv, ok := lv.(LBool); ok {
		return bool(lv)
	}

	L.RaiseError("must be bool , got %s", lv.Type().String())
	return false
}

func CheckTable(L *LState, lv LValue) *LTable {
	if lv, ok := lv.(*LTable); ok {
		return lv
	}
	L.RaiseError("must be LTable, got %s", lv.Type().String())
	return nil
}

func CheckFunction(L *LState, lv LValue) *LFunction {
	if lv, ok := lv.(*LFunction); ok {
		return lv
	}
	L.RaiseError("must be Function, got %s", lv.Type().String())
	return nil
}
func CheckSocket(v string) error {
	s := strings.Split(v, ":")
	if len(s) != 2 {
		return InvalidFormat
	}

	if net.ParseIP(s[0]) == nil {
		return InvalidIP
	}

	port, err := strconv.Atoi(s[1])
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return InvalidPort
	}
	return nil
}

func CheckIO(val *VelaData) IO {
	obj, ok := val.Data.(IO)
	if ok {
		return obj
	}
	return nil
}

func CheckWriter(val *VelaData) Writer {
	obj, ok := val.Data.(Writer)
	if ok {
		return obj
	}
	return nil
}

func CheckReader(val *VelaData) Reader {
	obj, ok := val.Data.(Reader)
	if ok {
		return obj
	}
	return nil
}

func CheckCloser(val *VelaData) Closer {
	obj, ok := val.Data.(Closer)
	if !ok {
		return nil
	}

	return obj
}

func S2B(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func S2L(s string) LString {
	return LString(s)
}

func B2L(b []byte) LString {
	return *(*LString)(unsafe.Pointer(&b))
}

func JsonMarshal(L *LState, i interface{}) LString {
	data, e := json.Marshal(i)
	if e != nil {
		L.RaiseError("%v", e)
		return LSNull
	}
	return *(*LString)(unsafe.Pointer(&data))
}

func L2SS(L *LState) []string {
	n := L.GetTop()
	if n == 0 {
		return nil
	}
	var ssv []string
	for i := 1; i <= n; i++ {
		lv := L.Get(i)
		if lv.Type() == LTNil {
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

func FileSuffix(path string) string {
	suffix := filepath.Ext(path)
	return strings.TrimSuffix(path, suffix)
}

func NewFunction(gn LGFunction) *LFunction {
	return &LFunction{
		IsG:       true,
		Proto:     nil,
		GFunction: gn,
	}
}

func CreateTable(acap, hcap int) *LTable {
	return newLTable(acap, hcap)
}

func CloneTable(v *LTable) *LTable {
	tab := &LTable{
		Metatable: v.Metatable,
	}

	a := len(v.array)
	if a > 0 {
		tab.array = make([]LValue, a)
		for i := 0; i < a; i++ {
			tab.array[i] = v.array[i]
		}
	}

	d := len(v.dict)
	if d > 0 {
		tab.dict = make(map[LValue]LValue, d)
		for key, val := range v.dict {
			tab.dict[key] = val
		}
	}
	s := len(v.strdict)
	if s > 0 {
		tab.strdict = make(map[string]LValue, d)
		for key, val := range v.strdict {
			tab.strdict[key] = val
		}
	}

	k := len(v.keys)
	if k > 0 {
		tab.array = make([]LValue, k)
		for i := 0; i < k; i++ {
			tab.keys[i] = v.keys[i]
		}
	}

	ki := len(v.k2i)
	if ki > 0 {
		tab.k2i = make(map[LValue]int, ki)
		for key, val := range v.k2i {
			tab.k2i[key] = val
		}
	}

	return tab
}
