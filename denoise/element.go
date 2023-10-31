package denoise

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/strutil"
)

const (
	String ElementType = iota + 1
	IndexEx
	FieldEx
	Noop
)

type ElementType uint8

type element struct {
	eType ElementType
	vm    *lua.LState
	Raw   []byte
	luaEx lua.IndexEx
	Field cond.FieldEx
	Value interface{}
}

func (e *element) v(key string) []byte {
	switch e.eType {
	case Noop:
		return nil
	case String:
		return e.Raw
	case IndexEx:
		return strutil.S2B(e.luaEx.Index(e.vm, key).String())
	case FieldEx:
		return strutil.S2B(e.Field.Field(key))
	default:
		return nil
	}
}
