package third

import (
	"encoding/xml"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"os"
)

func (i *Info) Cat() ([]byte, error) {
	if i.Err != nil {
		return nil, i.Err
	}

	chunk, err := os.ReadFile(i.Entry.File())
	return chunk, err
}

func (i *Info) BindJSON(L *lua.LState) int {
	chunk, err := i.Cat()
	if err != nil {
		L.Push(lua.LNil)
		L.Pushf("%v", err)
		return 2
	}

	tab, err := kind.Decode(L, chunk)
	if err != nil {
		L.Push(lua.LNil)
		L.Pushf("%v", err)
		return 2
	}

	L.Push(tab)
	return 1
}

func (i *Info) BindFastJSON(L *lua.LState) int {
	chunk, err := i.Cat()
	if err != nil {
		L.Push(lua.LNil)
		L.Pushf("%v", err)
		return 2
	}

	fast := &kind.Fast{}
	if e := fast.ParseBytes(chunk); e != nil {
		L.Push(lua.LNil)
		L.Pushf("%v", err)
		return 2
	}

	L.Push(fast)
	return 1
}

func (i *Info) BindXML(L *lua.LState) int {
	chunk, err := i.Cat()
	if err != nil {
		L.Push(lua.LNil)
		L.Pushf("%v", err)
		return 2
	}

	var v interface{}

	if e := xml.Unmarshal(chunk, &v); e != nil {
		L.Push(lua.LNil)
		L.Pushf("%v", e)
		return 2
	}

	tab := kind.DecodeValue(L, v)
	L.Push(tab)
	return 1
}

func (i *Info) BindLine(L *lua.LState) int {
	if i.Err != nil {
		L.Push(lua.LNil)
		L.Pushf("%v", i.Err)
		return 2
	}

	L.Push(&Scanner{file: i.Entry.File()})

	return 1
}
