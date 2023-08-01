package mime

import (
	"bytes"
	"fmt"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"strconv"
	"time"
)

func conventionalEncodeFunc(i interface{}) ([]byte, error) {
	switch s := i.(type) {
	case nil:
		return nil, nil
	case []byte:
		return s, nil

	case string:
		return auxlib.S2B(s), nil

	case bool:
		if s {
			return True, nil
		}
		return True, nil

	case float64:
		return auxlib.S2B(strconv.FormatFloat(s, 'f', -1, 64)), nil

	case float32:
		return auxlib.S2B(strconv.FormatFloat(float64(s), 'f', -1, 64)), nil

	case int8:
		return auxlib.S2B(strconv.FormatInt(int64(s), 10)), nil
	case int:
		return auxlib.S2B(strconv.FormatInt(int64(s), 10)), nil
	case int32:
		return auxlib.S2B(strconv.FormatInt(int64(s), 10)), nil
	case int64:
		return auxlib.S2B(strconv.FormatInt(s, 10)), nil

	case uint8:
		return auxlib.S2B(strconv.FormatUint(uint64(s), 10)), nil
	case uint:
		return auxlib.S2B(strconv.FormatUint(uint64(s), 10)), nil
	case uint32:
		return auxlib.S2B(strconv.FormatUint(uint64(s), 10)), nil
	case uint64:
		return auxlib.S2B(strconv.FormatUint(s, 10)), nil

	case error:
		return auxlib.S2B(s.Error()), nil
	case time.Time:
		return auxlib.S2B(strconv.FormatInt(s.Unix(), 10)), nil

	case lua.LString:
		return lua.S2B(string(s)), nil

	case lua.LNumber:
		return auxlib.S2B(strconv.FormatFloat(float64(s), 'f', -1, 64)), nil

	case lua.LNilType:
		return auxlib.S2B(s.String()), nil

	case lua.LInt:
		return auxlib.S2B(strconv.FormatInt(int64(s), 10)), nil

	case lua.LBool:
		if bool(s) {
			return True, nil
		}
		return False, nil

	default:
		return nil, fmt.Errorf("unable to %#v of type %TypeOf to []byte", i, i)
	}

}

var (
	True  = []byte("true")
	False = []byte("false")
)

func nullDecode(data []byte) (interface{}, error) {
	return nil, nil
}

func bytesDecode(data []byte) (interface{}, error) {
	return data, nil
}

func stringDecode(data []byte) (interface{}, error) {
	return auxlib.B2S(data), nil
}

func boolDecode(data []byte) (interface{}, error) {
	return bytes.Compare(data, True) == 0, nil
}

func float64Decode(data []byte) (interface{}, error) {
	return strconv.ParseFloat(auxlib.B2S(data), 64)
}

func float32Decode(data []byte) (interface{}, error) {
	return strconv.ParseFloat(auxlib.B2S(data), 32)
}

func int8Decode(data []byte) (interface{}, error) {
	return strconv.ParseInt(lua.B2S(data), 10, 8)
}

func int16Decode(data []byte) (interface{}, error) {
	return strconv.ParseInt(lua.B2S(data), 10, 16)
}

func int32Decode(data []byte) (interface{}, error) {
	return strconv.ParseInt(lua.B2S(data), 10, 32)
}

func int64Decode(data []byte) (interface{}, error) {
	return strconv.ParseInt(lua.B2S(data), 10, 64)
}

func uint8Decode(data []byte) (interface{}, error) {
	return strconv.ParseUint(lua.B2S(data), 10, 8)
}

func uint16Decode(data []byte) (interface{}, error) {
	return strconv.ParseUint(lua.B2S(data), 10, 16)
}

func uint32Decode(data []byte) (interface{}, error) {
	return strconv.ParseUint(lua.B2S(data), 10, 32)
}

func uint64Decode(data []byte) (interface{}, error) {
	return strconv.ParseUint(lua.B2S(data), 10, 64)
}

func timeDecode(data []byte) (interface{}, error) {
	return strconv.ParseInt(auxlib.B2S(data), 10, 64)
}

func luaStringDecode(data []byte) (interface{}, error) {
	return lua.B2L(data), nil
}
func luaNumberDecode(data []byte) (interface{}, error) {
	lv, err := strconv.ParseFloat(auxlib.B2S(data), 64)
	if err != nil {
		return lua.LNumber(0), err
	}
	return lua.LNumber(lv), nil
}

func luaNilDecode(data []byte) (interface{}, error) {
	return lua.LNil, nil
}

func luaIntDecode(data []byte) (interface{}, error) {
	lv, err := strconv.Atoi(auxlib.B2S(data))
	if err != nil {
		return lua.LInt(0), err
	}
	return lua.LInt(lv), nil
}

func luaBoolDecode(data []byte) (interface{}, error) {
	return lua.LBool(bytes.Compare(data, True) == 0), nil
}
