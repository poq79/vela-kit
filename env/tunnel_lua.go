package env

import (
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

func (env *Environment) sendL(L *lua.LState) int {

	n := L.GetTop()
	if n < 2 {
		L.Pushf("bad argument length, got %d", n)
		return 1
	}

	uri := L.IsString(1)
	if len(uri) == 0 {
		L.Pushf("bad argument #1 must be string")
		return 1
	}

	for i := 2; i <= n; i++ {
		item := L.Get(i)
		switch item.Type() {
		case lua.LTNil:
			continue
		case lua.LTString, lua.LTObject, lua.LTMap, lua.LTSlice, lua.LTVelaData, lua.LTAnyData:
			raw := item.String()
			if len(raw) == 0 {
				continue
			}

			err := env.Push(uri, json.RawMessage(raw))
			if err != nil {
				env.Errorf("push %s fail %v uri:%s", raw, err, uri)
				continue
			}
		}
	}

	return 0
}

/*

net.pipe("tcp://127.0.0.1:9092" , net.proxy("192.168.100.101:9092"))


*/

func (env *Environment) reconnectL(L *lua.LState) int {
	chains := pipe.NewByLua(L)
	env.onReconnect = chains
	return 0
}

func (env *Environment) tunnelIndexL(L *lua.LState, key string) lua.LValue {
	switch key {
	case "send":
		return lua.NewFunction(env.sendL)
	case "on_reconnect":
		return lua.NewFunction(env.reconnectL)
	}
	return lua.LNil
}
