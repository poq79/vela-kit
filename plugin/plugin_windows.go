package plugin

import (
	"github.com/vela-ssoc/vela-kit/env"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	xEnv vela.Environment
)

func environment() *env.Environment {
	v, ok := xEnv.(*env.Environment)
	if !ok {
		return nil
	}

	return v
}

func newLuaPlugin(L *lua.LState) int {
	path := L.CheckString(1)
	dll, err := syscall.LoadLibrary(path)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	inject, err := syscall.GetProcAddress(dll, "WithEnv")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	evi := environment()
	if evi == nil {
		L.RaiseError("invalid environment")
		return 0
	}

	_, _, e := syscall.Syscall(inject, 1, uintptr(unsafe.Pointer(&evi)), 0, 0)
	if int(e) == 0 {
		L.RaiseError("%s", e.Error())
	}

	return 0
}

func Constructor(env vela.Environment) {
	xEnv = env
	env.Infof("plugin running in %s", runtime.GOOS)
	env.Set("plugin", lua.NewFunction(newLuaPlugin))
}
