package env

import (
	"github.com/shirou/gopsutil/host"
	"github.com/vela-ssoc/vela-kit/lua"
	"net"
)

func (env *Environment) nodeIDL(L *lua.LState) int {
	L.Push(lua.S2L(env.ID()))
	return 1
}

func (env *Environment) Kernel() string {
	info, err := host.Info()
	if err != nil {
		env.log.Error("search linux kernel info fail %v", err)
		return ""
	}

	return info.KernelVersion
}

func (env *Environment) inetL(L *lua.LState) int {
	L.Push(lua.S2L(env.Inet()))
	return 1
}

func (env *Environment) inet6L(L *lua.LState) int {
	L.Push(lua.LSNull)
	return 0
}

func (env *Environment) macL(L *lua.LState) int {
	L.Push(lua.S2L(env.Mac()))
	return 1
}

func (env *Environment) archL(L *lua.LState) int {
	L.Push(lua.S2L(env.Arch()))
	return 1
}

func (env *Environment) addrL(L *lua.LState) int {
	v := lua.Slice{}
	ifat, err := net.InterfaceAddrs()
	if err != nil {
		env.Errorf("not found interface %v", err)
		L.Push(v)
		return 1
	}

	for _, addr := range ifat {
		ip, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if ip.IP.IsLoopback() {
			continue
		}

		v = append(v, lua.S2L(ip.IP.String()))
	}
	L.Push(v)
	return 1
}
