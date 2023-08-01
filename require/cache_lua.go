package require

import (
	"github.com/vela-ssoc/vela-kit/lua"
)

func (c *cache) Type() lua.LValueType {
	return lua.LTObject
}

func (c *cache) gc() {
	defer func() {
		c.status = DEL
		c.cdata = lua.LNil
	}()

	var gc lua.LValue

	switch v := c.cdata.(type) {
	case lua.IndexEx:
		gc = v.Index(c.co, "_gc")
	case *lua.LTable:
		gc = v.RawGetString("_gc")
	default:
		return
	}

	switch gc.Type() {
	case lua.LTNil:
		return

	case lua.LTFunction:
		cp := xEnv.P(gc.(*lua.LFunction))
		c.co.CallByParam(cp)
	default:
		xEnv.Errorf("got gc field , but not function")
	}
}

func (c *cache) object() lua.LValue {
	if c.err != nil {
		return lua.LNil
	}

	if c.status != DEL {
		return c.cdata
	}

	return instance.require(c.co, c.name)
}

func (c *cache) String() string {
	return c.object().String()
}

func (c *cache) AssertFloat64() (float64, bool) {
	return c.object().AssertFloat64()
}

func (c *cache) AssertString() (string, bool) {
	return c.object().AssertString()
}

func (c *cache) AssertFunction() (*lua.LFunction, bool) {
	return c.object().AssertFunction()
}

func (c *cache) Peek() lua.LValue {
	return c.object()
}
