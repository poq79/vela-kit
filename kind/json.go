package kind

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
)

var empty = lua.S2L("[]")

var (
	errNested      = errors.New("cannot encode recursively nested tables to JSON")
	errSparseArray = errors.New("cannot encode sparse array")
	errInvalidKeys = errors.New("cannot encode mixed or invalid key types")
)

type invalidTypeError lua.LValueType

func (i invalidTypeError) Error() string {
	return `cannot encode ` + lua.LValueType(i).String() + ` to JSON`
}

// Encode returns the JSON encoding of value.
func Encode(value lua.LValue) ([]byte, error) {
	return json.Marshal(jsonValue{
		LValue:  value,
		visited: make(map[*lua.LTable]bool),
	})
}

type jsonValue struct {
	lua.LValue
	visited map[*lua.LTable]bool
}

func (j jsonValue) MarshalJSON() (data []byte, err error) {
	switch converted := j.LValue.(type) {
	case lua.LBool:
		data, err = json.Marshal(bool(converted))
	case lua.LNumber:
		data, err = json.Marshal(float64(converted))
	case lua.LInt:
		data, err = json.Marshal(int(converted))
	case *lua.LNilType:
		data = []byte(`null`)
	case lua.LString:
		data, err = json.Marshal(string(converted))
	case *lua.LTable:
		if j.visited[converted] {
			return nil, errNested
		}
		j.visited[converted] = true

		key, value := converted.Next(lua.LNil)

		switch key.Type() {
		case lua.LTNil: // empty table
			data = []byte(`[]`)
		case lua.LTNumber:
			arr := make([]jsonValue, 0, converted.Len())
			expectedKey := lua.LNumber(1)
			for key != lua.LNil {
				if key.Type() != lua.LTNumber {
					err = errInvalidKeys
					return
				}
				if expectedKey != key {
					err = errSparseArray
					return
				}
				arr = append(arr, jsonValue{value, j.visited})
				expectedKey++
				key, value = converted.Next(key)
			}
			data, err = json.Marshal(arr)
		case lua.LTString:
			obj := make(map[string]jsonValue)
			for key != lua.LNil {
				if key.Type() != lua.LTString {
					err = errInvalidKeys
					return
				}
				obj[key.String()] = jsonValue{value, j.visited}
				key, value = converted.Next(key)
			}
			data, err = json.Marshal(obj)
		default:
			err = errInvalidKeys
		}
	default:
		err = invalidTypeError(j.LValue.Type())
	}
	return
}

// Decode converts the JSON encoded data to Lua values.
func Decode(L *lua.LState, data []byte) (lua.LValue, error) {
	var value interface{}
	err := json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}
	return DecodeValue(L, value), nil
}

// DecodeValue converts the value to a Lua value.
//
// This function only converts values that the encoding/json package decodes to.
// All other values will return lua.LNil.
func DecodeValue(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case int:
		return lua.LInt(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case json.Number:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(DecodeValue(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), DecodeValue(L, item))
		}
		return tbl
	case nil:
		return lua.LNil
	}

	return lua.LNil
}

func Arr2Json(arr []lua.LValue, enc *JsonEncoder) (err error) {
	enc.Arr("")
	for _, lv := range arr {
		err = LValueEncode(lv, enc)
	}
	enc.End("],")
	return
}

func Tab2Json(tab *lua.LTable, enc *JsonEncoder) (err error) {
	if tab.IsArray() {
		return Arr2Json(tab.Array(), enc)
	}

	enc.Tab("")
	tab.RangeStrDict(func(key string, lv lua.LValue) bool {
		enc.Key(key)
		err = LValueEncode(lv, enc)
		return err == nil
	})

	if err != nil {
		return
	}

	tab.RangeDict(func(key lua.LValue, val lua.LValue) bool {
		switch key.Type() {
		case lua.LTNil:
			//todo
		case lua.LTString, lua.LTInt, lua.LTNumber:
			enc.Key(val.String())
			err = LValueEncode(val, enc)
		default:
			err = fmt.Errorf("invalid table field type , got %s", key.Type().String())
		}

		return err == nil

	})
	enc.End("},")
	return
}

func Object2Json(data interface{}, enc *JsonEncoder) error {

	switch item := data.(type) {
	case nil:
		//todo
	case interface{ Json() []byte }:
		enc.Write(item.Json())
	case interface{ Byte() []byte }:
		enc.Write(item.Byte())
	case interface{ String() string }:
		enc.WriteString(item.String())

	default:
		chunk, err := json.Marshal(data)
		if err != nil {
			return err
		}
		enc.Write(chunk)
	}
	return nil
}

func LValueEncode(lv lua.LValue, enc *JsonEncoder) (err error) {
	switch lv.Type() {
	case lua.LTString:
		enc.Val(lv.String())
		enc.WriteByte(',')

	case lua.LTTable:
		err = Tab2Json(lv.(*lua.LTable), enc)

	case lua.LTAnyData:
		err = Object2Json(lv.(*lua.AnyData).Data, enc)

	case lua.LTUserData:
		err = Object2Json(lv.(*lua.LUserData).Value, enc)

	case lua.LTVelaData:
		err = Object2Json(lv.(*lua.VelaData).Data, enc)
	case lua.LTObject:
		err = Object2Json(lv, enc)

	default:
		//LBool , LNumber
		enc.WriteString(lv.String())
		enc.WriteByte(',')
	}

	return nil
}

func Marshal(lv lua.LValue) ([]byte, error) {
	enc := NewJsonEncoder()
	err := LValueEncode(lv, enc)
	if err != nil {
		return nil, err
	}

	enc.End("") //取消最后一个逗号

	return enc.Bytes(), nil
}
