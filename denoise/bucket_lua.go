package denoise

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

func (bkt *Bucket) String() string                         { return "" }
func (bkt *Bucket) Type() lua.LValueType                   { return lua.LTObject }
func (bkt *Bucket) AssertFloat64() (float64, bool)         { return 0, false }
func (bkt *Bucket) AssertString() (string, bool)           { return "", false }
func (bkt *Bucket) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (bkt *Bucket) Peek() lua.LValue                       { return bkt }

func (bkt *Bucket) SectionL(L *lua.LState) int {
	tab := L.CheckTable(1)
	size := auxlib.ToInt(tab.RawGetString("size").String())
	if size == 0 {
		size = 8192
	}
	s := NewSection(size)

	tab.Range(func(key string, val lua.LValue) {
		s.NewIndex(L, key, val)
	})

	bkt.Add(s)
	return 0
}

func (bkt *Bucket) ignoreL(L *lua.LState) int {
	return 0
}

func (bkt *Bucket) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "section":
		return lua.NewFunction(bkt.SectionL)

	case "ignore":
		return lua.NewFunction(bkt.ignoreL)
	default:
		return lua.LNil
	}
}

func NewBucketL(L *lua.LState) *Bucket {
	bkt := &Bucket{
		co:     xEnv.Clone(L),
		ignore: cond.NewIgnore(),
	}
	return bkt
}
