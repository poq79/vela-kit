package mime

import (
	"errors"
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
	"reflect"
	"sync"
	"time"
)

var (
	xEnv  vela.Environment
	mutex sync.RWMutex
)

var (
	//DefaultEncode = BinaryEncode
	//DefaultDecode = BinaryDecode

	//mimeEncode = make(map[string]assert.EncodeFunc)
	//mimeDecode = make(map[string]assert.DecodeFunc)

	notFoundEncode = errors.New("not found mime encode")
)

func Encode(v interface{}) ([]byte, string, error) {
	return xEnv.MimeEncode(v)
}

func Decode(name string, data []byte) (interface{}, error) {
	return xEnv.MimeDecode(name, data)
}

func Name(v interface{}) string {
	if v == nil {
		return "nil"
	}

	vt := reflect.TypeOf(v)

LOOP:
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
		goto LOOP
	}
	return vt.String()
}

func Register(v interface{}, encode vela.EncodeFunc, decode vela.DecodeFunc) {
	xEnv.Mime(v, encode, decode)
}

func Constructor(env vela.Environment) {
	xEnv = env
	env.Mime(nil, conventionalEncodeFunc, nullDecode)
	env.Mime("", conventionalEncodeFunc, stringDecode)
	env.Mime([]byte{}, conventionalEncodeFunc, bytesDecode)
	env.Mime(true, conventionalEncodeFunc, boolDecode)
	env.Mime(float64(0), conventionalEncodeFunc, float64Decode)
	env.Mime(float32(0), conventionalEncodeFunc, float32Decode)
	env.Mime(int(0), conventionalEncodeFunc, int32Decode)
	env.Mime(int8(0), conventionalEncodeFunc, int8Decode)
	env.Mime(int16(0), conventionalEncodeFunc, int16Decode)
	env.Mime(int32(0), conventionalEncodeFunc, int32Decode)
	env.Mime(int64(0), conventionalEncodeFunc, int64Decode)
	env.Mime(uint(0), conventionalEncodeFunc, uint32Decode)
	env.Mime(uint8(0), conventionalEncodeFunc, uint8Decode)
	env.Mime(uint16(0), conventionalEncodeFunc, uint16Decode)
	env.Mime(uint32(0), conventionalEncodeFunc, uint32Decode)
	env.Mime(uint64(0), conventionalEncodeFunc, uint64Decode)
	env.Mime(time.Now(), conventionalEncodeFunc, timeDecode)
	env.Mime(lua.LString(""), conventionalEncodeFunc, luaStringDecode)
	env.Mime(lua.LBool(true), conventionalEncodeFunc, luaBoolDecode)
	env.Mime(lua.LNumber(0), conventionalEncodeFunc, luaNumberDecode)
	env.Mime(lua.LInt(0), conventionalEncodeFunc, luaIntDecode)
	env.Mime(lua.LNil, conventionalEncodeFunc, luaNilDecode)
}
