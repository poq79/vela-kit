package lua

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

/*
	type A struct {
		name string `lua:"name"`
		Age  int  `lua:"age"`
	}

	func(a *A) Show1(L *LState) int {}
	func(a *A) Show2(L *LState) LValue {}
	func(a *A) LuaShow() string {}

	local a = rock.A{}

	print(a.name)
	print(a.age)
	a.show1()
	a.show2()
	local r = a.show3()
*/

type void struct {
	t *struct{}
	p unsafe.Pointer
}

type AnyData struct {
	Data interface{}
	meta UserKV

	flag     int
	ptr      uintptr
	metaFunc UserKV
	metaElem map[string]reflect.StructField
}

func (ad *AnyData) Len(L *LState) int {
	switch data := ad.Data.(type) {
	case string:
		return len(data)
	case []byte:
		return len(data)
	case interface{ Len() int }:
		return data.Len()
	default:
		L.RaiseError("cant not length")
		return 0
	}
}

type index interface {
	Index(*LState, string) LValue
}

type newIndex interface {
	NewIndex(*LState, string, LValue)
}

//NewAnyData(v , lua.Reflect(true)) {}

const (
	OFF int = iota + 1
	ALL
	ELEM
	FUNC
)

func Reflect(flag int) func(*AnyData) {
	return func(ad *AnyData) {
		ad.flag = flag

		switch flag {
		case OFF:
			//todo

		case ELEM:
			ad.ptr = uintptr((*void)(unsafe.Pointer(&ad.Data)).p)
			ad.initMetaElem()

		case FUNC:
			ad.ptr = uintptr((*void)(unsafe.Pointer(&ad.Data)).p)
			ad.initMetaFunc()

		case ALL:
			ad.initMetaElem()
			ad.initMetaFunc()
			ad.ptr = uintptr((*void)(unsafe.Pointer(&ad.Data)).p)
		}

	}
}

func NewAnyData(v interface{}, opts ...func(data *AnyData)) *AnyData {
	ad := &AnyData{Data: v, flag: OFF}
	for _, fn := range opts {
		fn(ad)
	}

	return ad
}

func (ad *AnyData) Type() LValueType                   { return LTAnyData }
func (ad *AnyData) AssertFloat64() (float64, bool)     { return 0, false }
func (ad *AnyData) AssertString() (string, bool)       { return "", false }
func (ad *AnyData) AssertFunction() (*LFunction, bool) { return nil, false }
func (ad *AnyData) Peek() LValue                       { return ad }

func (ad *AnyData) String() string {
	switch val := ad.Data.(type) {

	case interface{ String() string }:
		return val.String()

	case interface{ Byte() []byte }:
		return B2S(val.Byte())

	default:
		return fmt.Sprintf("AnyData: %p", ad)
	}
}

func (ad *AnyData) Meta(key string, val LValue) {
	if val == LNil {
		return
	}

	if ad.meta == nil {
		ad.meta = NewUserKV()
	}

	ad.meta.Set(key, val)
}

func (ad *AnyData) initMetaFunc() {
	//初始化方法
	ad.metaFunc = NewUserKV()

	rt := reflect.TypeOf(ad.Data)
	n := rt.NumMethod()
	if n == 0 {
		return
	}

	rv := reflect.ValueOf(ad.Data)
	for i := 0; i < n; i++ {
		m := rt.Method(i)
		vi := rv.Method(i).Interface()
		//选择
		switch fn := vi.(type) {

		case func():
			ad.metaFunc.Set(m.Name, newLFunctionG(func(_ *LState) int {
				fn()
				return 0
			}, nil, 0))

		case func() int:
			ad.metaFunc.Set(m.Name, newLFunctionG(func(_ *LState) int {
				fn()
				return 0
			}, nil, 0))

		case func() []byte:
			ad.metaFunc.Set(m.Name, newLFunctionG(func(co *LState) int {
				data := fn()
				co.Push(B2L(data))
				return 1
			}, nil, 0))

		case func() string:
			ad.metaFunc.Set(m.Name, newLFunctionG(func(co *LState) int {
				data := fn()
				co.Push(S2L(data))
				return 1
			}, nil, 0))

		case func(*LState) int:
			ad.metaFunc.Set(m.Name, newLFunctionG(fn, nil, 0))

		case func(*LState) LValue:
			ad.metaFunc.Set(m.Name, newLFunctionG(func(co *LState) int {
				co.Push(fn(co))
				return 1
			}, nil, 0))
		}
	}

}

func (ad *AnyData) initMetaElem() {

	ad.metaElem = make(map[string]reflect.StructField, 5)

	rt := reflect.TypeOf(ad.Data)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	nE := rt.NumField()
	for i := 0; i < nE; i++ {
		t := rt.Field(i)
		tag, ok := t.Tag.Lookup("lua")
		if ok {
			if tag == "_" {
				continue
			}
			tag = strings.Split(tag, ",")[0]
		} else {
			tag = t.Name
		}

		switch t.Type.Kind() {
		case reflect.String, reflect.Bool, reflect.Interface,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64,
			reflect.Float32, reflect.Float64:
			ad.metaElem[tag] = t

		case reflect.Struct:
			switch t.Type.String() {
			case "time.Time":
				ad.metaElem[tag] = t
			}

		default:
			continue
		}
	}
}

func (ad *AnyData) getElem(L *LState, key string) LValue {

	if ad.flag == FUNC || ad.flag == OFF {
		return LNil
	}

	field, ok := ad.metaElem[key]
	if !ok {
		return LNil
	}

	switch field.Type.Kind() {
	case reflect.String:
		return LString(*(*string)(unsafe.Pointer(ad.ptr + field.Offset)))
	case reflect.Int:
		return LNumber(*(*int)(unsafe.Pointer(ad.ptr + field.Offset)))
	case reflect.Uint8:
		return LNumber(*(*uint8)(unsafe.Pointer(ad.ptr + field.Offset)))
	case reflect.Uint16:
		return LNumber(*(*uint16)(unsafe.Pointer(ad.ptr + field.Offset)))
	case reflect.Uint32:
		return LNumber(*(*uint32)(unsafe.Pointer(ad.ptr + field.Offset)))
	case reflect.Uint64:
		return LNumber(*(*uint64)(unsafe.Pointer(ad.ptr + field.Offset)))

	case reflect.Int8:
		return LNumber(*(*int8)(unsafe.Pointer(ad.ptr + field.Offset)))
	case reflect.Int16:
		return LNumber(*(*int16)(unsafe.Pointer(ad.ptr + field.Offset)))
	case reflect.Int32:
		return LNumber(*(*int32)(unsafe.Pointer(ad.ptr + field.Offset)))
	case reflect.Int64:
		return LNumber(*(*int64)(unsafe.Pointer(ad.ptr + field.Offset)))

	case reflect.Bool:
		return LBool(*(*bool)(unsafe.Pointer(ad.ptr + field.Offset)))

	case reflect.Interface:
		return L.NewAnyData(*(*interface{})(unsafe.Pointer(ad.ptr + field.Offset)))

	case reflect.Struct:
		switch field.Type.String() {

		case "time.Time":
			item := *(*time.Time)(unsafe.Pointer(ad.ptr + field.Offset))
			return S2L(item.Format(time.RFC3339Nano))

		default:
			return LNil

		}

	default:
		return LNil
	}

}

func (ad *AnyData) setElem(L *LState, key string, val LValue) {
	field, ok := ad.metaElem[key]
	if !ok {
		return
	}

	rv := reflect.ValueOf(ad.Data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	vt := val.Type()

	switch field.Type.Kind() {
	case reflect.String:
		if vt == LTString {
			*(*string)(unsafe.Pointer(ad.ptr + field.Offset)) = val.String()
			return
		}
		L.RaiseError("%s args must be string , got %s", key, val.Type().String())

	case reflect.Int:
		if vt == LTNumber {
			*(*int)(unsafe.Pointer(ad.ptr + field.Offset)) = int(LVAsNumber(val))
			return
		}
		L.RaiseError("%s args must be number, got %s", key, val.Type().String())

	case reflect.Bool:
		if vt == LTBool {
			*(*bool)(unsafe.Pointer(ad.ptr + field.Offset)) = bool(val.(LBool))
			return
		}
		L.RaiseError("%s args must be bool , got %s", key, val.Type().String())

	default:
		L.RaiseError("can't set %s elem", key)
	}

}

func (ad *AnyData) inMeta(key string) (LValue, bool) {
	if ad.meta == nil {
		return nil, false
	}

	return ad.meta.V(key)
}

func (ad *AnyData) Getter(L *LState, key string) (LValue, bool) {
	obj, ok := ad.Data.(index)
	if !ok {
		return nil, false
	}

	lv := obj.Index(L, key)
	if lv == nil {
		return nil, false
	}
	if lv == LNil {
		return nil, false
	}

	return lv, true
}

func (ad *AnyData) Index(L *LState, key string) LValue {
	if lv, ok := ad.inMeta(key); ok {
		return lv
	}

	if lv, ok := ad.Getter(L, key); ok {
		return lv
	}

	switch ad.flag {
	case OFF:
		return LNil

	case ELEM:
		return ad.getElem(L, key)
	case FUNC:
		return ad.metaFunc.Get(key)
	case ALL:
		if lv := ad.getElem(L, key); lv != nil {
			return lv
		}
		return ad.metaFunc.Get(key)

	default:
		return LNil
	}
}

func (ad *AnyData) NewIndex(L *LState, key string, val LValue) {

	obj, ok := ad.Data.(newIndex)
	if ok {
		obj.NewIndex(L, key, val)
		return
	}

	switch ad.flag {
	case ALL, ELEM:
		ad.setElem(L, key, val)
	}
}

func (ls *LState) NewAnyData(v interface{}, opts ...func(*AnyData)) *AnyData {
	return NewAnyData(v, opts...)
}

func (ls *LState) CheckAnyData(n int) *AnyData {
	v := ls.Get(n)
	if lv, ok := v.(*AnyData); ok {
		return lv
	}
	ls.TypeError(n, LTAnyData)
	return nil
}
