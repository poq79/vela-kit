package hashmap

import (
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/reflectx"
)

var hmMetaTab = map[string]lua.LValue{
	"range": lua.NewFunction(hmMetaRange),
}

func (hm HMap) String() string                         { return lua.B2S(hm.Byte()) }
func (hm HMap) Type() lua.LValueType                   { return lua.LTObject }
func (hm HMap) AssertFloat64() (float64, bool)         { return 0, false }
func (hm HMap) AssertString() (string, bool)           { return "", false }
func (hm HMap) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (hm HMap) Peek() lua.LValue                       { return hm }

func (hm HMap) Byte() []byte {
	chunk, err := json.Marshal(hm)
	if err != nil {
		xEnv.Errorf("hashmap to bytes fail %v", err)
		return nil
	}

	return chunk
}

func (hm HMap) MetaTable(L *lua.LState, key string) lua.LValue {
	lv, ok := hmMetaTab[key]
	if ok {
		return lv
	}
	return lua.LNil
}

func (hm HMap) NewMeta(L *lua.LState, key lua.LValue, val lua.LValue) {
	if key.Type() != lua.LTString {
		L.RaiseError("lua table meta invalid type got %s , must string", key.Type().String())
		return
	}

	hm[key.String()] = val
}

func (hm HMap) Index(L *lua.LState, key string) lua.LValue {
	val, ok := hm[key]
	if ok {
		return reflectx.ToLValue(val, L)
	}
	return lua.LNil
}

func (hm HMap) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch val.Type() {
	case lua.LTNil:
		hm[key] = nil
		delete(hm, key)
	case lua.LTNumber:
		hm[key] = float64(val.(lua.LNumber))
	case lua.LTInt:
		hm[key] = int(val.(lua.LInt))
	case lua.LTBool:
		hm[key] = bool(val.(lua.LBool))
	default:
		hm[key] = val.String()
	}
}

func CheckHMap(L *lua.LState, idx int) HMap {
	obj := L.CheckObject(idx)

	if hm, ok := obj.(HMap); ok {
		return hm
	}
	L.RaiseError("invalid hashmap")
	return nil
}

func hmMetaRange(L *lua.LState) int {
	hm := CheckHMap(L, 1)
	cp := xEnv.P(L.CheckFunction(2))
	co := xEnv.Clone(L)
	defer xEnv.Free(co)

	for k, v := range hm {
		err := co.CallByParam(cp, lua.S2L(k), reflectx.ToLValue(v, L))
		if err != nil {
			xEnv.Errorf("hashmap range error %v", err)
			return 0
		}

		if co.IsTrue(-1) {
			return 0
		}
		co.SetTop(0)
	}

	return 0
}
