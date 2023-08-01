package proxy

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

/*
	vela.proxy("tcp://192.168.1.1:9090")

	p := struct {
		url: tcp://192.168.1.1:9090
		env: xEnv,
	}

	conn , err := p.dail()
	if err != nil {

	}


*/

func proxyL(L *lua.LState) int {
	addr := L.CheckString(1)
	L.Push(&Proxy{
		addr: addr,
	})
	return 1
}
func IsProxy(L *lua.LState, val lua.LValue) *Proxy {
	if val.Type() != lua.LTObject {
		L.RaiseError(" not proxy data got %s", val.Type().String())
		return nil
	}

	p, ok := val.(*Proxy)
	if !ok {
		L.RaiseError(" vela data not proxy")
		return nil
	}
	return p
}

func New(addr string) *Proxy {
	return &Proxy{addr: addr}
}

func Check(L *lua.LState, idx int) *Proxy {
	obj := L.CheckObject(idx)
	if obj == nil {
		return nil
	}

	p, ok := obj.(*Proxy)
	if ok {
		return p
	}

	L.RaiseError("#%d vela data not proxy", idx)
	return nil
}

func WithEnv(env vela.Environment) {
	xEnv = env
	xEnv.Set("proxy", lua.NewExport("net.proxy.export", lua.WithFunc(proxyL)))
}
