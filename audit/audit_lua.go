package audit

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

func (adt *Audit) toL(L *lua.LState) int {
	adt.cfg.sdk = auxlib.CheckWriter(L.CheckVelaData(1), L)
	return 0
}

func (adt *Audit) pipeL(L *lua.LState) int {
	adt.cfg.pipe.CheckMany(L, pipe.Seek(0))
	return 0
}

func (adt *Audit) passL(L *lua.LState) int {
	key := L.CheckString(1)
	filter := L.CheckString(2)
	adt.cfg.pass = append(adt.cfg.pass, newFilter(key, filter))
	return 0
}

func (adt *Audit) inhibitL(L *lua.LState) int {
	tag := L.CheckString(1)
	ttl := L.CheckInt(2)
	adt.cfg.rate = append(adt.cfg.rate, newInhibitMatch(tag, ttl))
	return 0
}

/*
	adt.init{}
	adt.pass()
	adt.pipe(_(ev) end)
	adt.to(sdk)
	adt.inhibit("$id.$inet.$from.$remote_addr.$subject" , 5 * 60) //5分钟 告警一次
*/

func (adt *Audit) Index(L *lua.LState, key string) lua.LValue {
	if !checkout(L) {
		return lua.LNil
	}

	switch key {

	case "pass":
		return lua.NewFunction(adt.passL)

	case "pipe":
		return lua.NewFunction(adt.pipeL)

	case "to":
		return lua.NewFunction(adt.toL)

	case "inhibit":
		return lua.NewFunction(adt.inhibitL)

	case "start":
		return lua.NewFunction(func(co *lua.LState) int {
			xEnv.Start(L, adt).From(co.CodeVM()).Do()
			return 0
		})

	default:

		//todo
		return lua.LNil
	}
}
