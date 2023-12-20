package mime

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
)

func LuaTableEncodeFunc(v interface{}) ([]byte, error) {
	tab, ok := v.(*lua.LTable)
	if ok {
		return kind.Encode(tab)
	}

	return nil, fmt.Errorf("not table type")
}

func LuaTableDecodeFunc(chunk []byte) (interface{}, error) {
	lv, err := kind.Decode(nil, chunk)
	return lv, err
}
