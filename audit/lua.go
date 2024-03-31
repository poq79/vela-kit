package audit

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

func checkout(L *lua.LState) bool {
	if !L.CheckCodeVM("audit") {
		L.RaiseError("audit not allow with %s", L.CodeVM())
		return false
	}

	return true
}

/*
	local adt = audit.new{
		file = "xxx"
	}

	adt.to(lua.writer)

	adt.pass("id" , "*helo")

	adt.pipe(_(ev) {
	})

	adt.pipe(service.a.kfk)

	adt.start()
*/

func (adt *Audit) NewL(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewVela(adt.Name(), typeof)
	if proc.IsNil() {
		adt.cfg = cfg
		proc.Set(adt)
	} else {
		adt.cfg = cfg
	}

	L.Push(proc)
	return 0
}

func Constructor(env vela.Environment, callback func(*Audit) error) {
	xEnv = env
	adt := New()
	xEnv.Set("adt", lua.NewExport("lua.audit.export", lua.WithFunc(adt.NewL)))
	kv := lua.NewUserKV()
	xEnv.Set("adt", kv)

	xEnv.Set("event", lua.NewFunction(newLuaEvent))
	xEnv.Set("Debug", lua.NewFunction(newLuaDebug))
	xEnv.Set("debug", lua.NewFunction(newLuaDebug))
	xEnv.Set("Error", lua.NewFunction(newLuaError))
	xEnv.Set("ERR", lua.NewFunction(newLuaError))
	xEnv.Set("Info", lua.NewFunction(newLuaInfo))
	xEnv.Set("T", lua.NewFunction(newLuaObjectType))

	if err := callback(adt); err != nil {
		xEnv.Errorf("audit callback fail %v", err)
	}

	adt.define(xEnv.R())
}
