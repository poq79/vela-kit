package kind

import (
	"github.com/valyala/fastjson"
	"github.com/valyala/fastjson/fastfloat"
	"github.com/vela-ssoc/vela-kit/lua"
)

type Fast struct {
	value *fastjson.Value
}

func (f *Fast) String() string                         { return f.value.String() }
func (f *Fast) Type() lua.LValueType                   { return lua.LTObject }
func (f *Fast) AssertFloat64() (float64, bool)         { return 0, false }
func (f *Fast) AssertString() (string, bool)           { return "", false }
func (f *Fast) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (f *Fast) Peek() lua.LValue                       { return f }

func (f *Fast) Int(L *lua.LState) int {
	key := L.CheckString(1)
	n := f.value.GetInt(key)
	L.Push(lua.LNumber(n))
	return 1
}

func (f *Fast) Str(L *lua.LState) int {
	key := L.CheckString(1)
	b := f.value.GetStringBytes(key)
	L.Push(lua.LString(lua.B2S(b)))
	return 1
}

func (f *Fast) Bool(L *lua.LState) int {
	key := L.CheckString(1)
	b := f.value.GetBool(key)
	L.Push(lua.LBool(b))
	return 1
}

func (f *Fast) ParseBytes(body []byte) error {
	v, err := fastjson.ParseBytes(body)
	if err != nil {
		return err
	}
	f.value = v

	return nil
}

func (f *Fast) Parse(body string) error {
	v, err := fastjson.Parse(body)
	if err != nil {
		return err
	}
	f.value = v
	return nil
}

func (f *Fast) visit(key string) lua.LValue {
	if f.value == nil {
		return lua.LNil
	}

	v := f.value.Get(key)
	if v == nil {
		return lua.LNil
	}

	switch v.Type() {
	case fastjson.TypeNull:
		return lua.LNil

	case fastjson.TypeString:
		item := v.String()
		return lua.S2L(item[1 : len(item)-1])

	case fastjson.TypeNumber:
		n, err := fastfloat.Parse(v.String())
		if err != nil {
			return lua.LNil
		}
		return lua.LNumber(n)

	case fastjson.TypeObject:
		return &Fast{value: v}

	case fastjson.TypeArray:
		return &Fast{value: v}

	case fastjson.TypeTrue:
		return lua.LTrue
	case fastjson.TypeFalse:
		return lua.LFalse

	default:
		return lua.S2L(v.String()) //typeRawString 7
	}

}

func (f *Fast) Index(L *lua.LState, key string) lua.LValue {
	return f.visit(key)
}

func (f *Fast) Meta(L *lua.LState, key lua.LValue) lua.LValue {
	return f.visit(key.String())
}

func (f *Fast) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch val.Type() {
	case lua.LTNil:
		return
	case lua.LTNumber, lua.LTBool, lua.LTInt, lua.LTInt64, lua.LTUint, lua.LTUint64:
		fv, err := fastjson.Parse(val.String())
		if err != nil {
			L.RaiseError("fastjson decode fail %v", err)
			return
		}

		f.value.Set(key, fv)
	case lua.LTObject:
		item, ok := val.(*Fast)
		if ok {
			f.value.Set(key, item.value)
			return
		}

		fv, err := fastjson.Parse(val.String())
		if err != nil {
			L.RaiseError("fastjson decode fail %v", err)
			return
		}

		f.value.Set(key, fv)
	case lua.LTString:
		f.value.Set(key, fastjson.MustParse("\""+val.String()+"\""))
	default:
		f.value.Set(key, fastjson.MustParse("\""+val.String()+"\""))
	}
}
