package env

import (
	"bytes"
	strutil "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/tasktree"
	vela "github.com/vela-ssoc/vela-kit/vela"
	"net"
	"strings"
	"time"
)

type Once struct {
	Time time.Time `json:"time" lua:"time"`
	Key  string    `json:"key" lua:"key"`
}

func (env *Environment) OnceL(L *lua.LState) int {
	key := L.CheckString(1)
	if e := strutil.Name(key); e != nil {
		L.RaiseError("once %v", e)
		return 0
	}

	bkt := env.Shm("VELA-EXDATA-DB", "ONCE-BKT")
	if bkt == nil {
		L.RaiseError("%v", "not found [VEAL-RUNTIME-DB,ONCE-BKT]")
		return 0
	}

	handle := func(ov *Once) int {
		ud := lua.NewAnyData(ov, lua.Reflect(lua.ELEM))
		chain := pipe.NewByLua(L, pipe.Seek(1), pipe.Env(env))
		chain.Do(ud, L, func(err error) {
			L.RaiseError("%v", err)
		})
		_ = bkt.Store(key, true, 0)
		L.Push(ud)
		return 1
	}

	val, err := bkt.Get(key)
	if err != nil {
		return handle(&Once{time.Now(), key})
	}

	ov, ok := val.(*Once)
	if !ok {
		return handle(&Once{time.Now(), key})
	}

	L.Push(lua.NewAnyData(ov, lua.Reflect(lua.ELEM)))
	return 1
}

func (env *Environment) hideIndexL(L *lua.LState, key string) lua.LValue {
	switch key {
	case "tag":
		s := &lua.Slice{}
		s.Build(env.hide.Tags)
		return s
	case "version":
		return lua.S2L(env.hide.Edition)
	case "servername":
		return lua.S2L(env.hide.Servername)
	default:
		return lua.LNil
	}
}

func (env *Environment) brokerIndexL(L *lua.LState, key string) lua.LValue {
	switch key {
	case "addr":
		if env.tnl == nil {
			return lua.LSNull
		}

		addr := env.tnl.RemoteAddr().(*net.TCPAddr)
		return lua.S2L(addr.IP.String())

	case "port":
		if env.tnl == nil {
			return lua.LInt(0)
		}
		addr := env.tnl.RemoteAddr().(*net.TCPAddr)
		return lua.LInt(addr.Port)

	default:
		return lua.LNil
	}
}

func (env *Environment) exdataIndexL(L *lua.LState, key string) lua.LValue {
	switch key {
	case "A":
		return lua.ToLValue(L.Metadata(0))
	case "B":
		return lua.ToLValue(L.Metadata(1))
	case "C":
		return lua.ToLValue(L.Metadata(2))
	case "D":
		return lua.ToLValue(L.Metadata(3))
	case "E":
		return lua.ToLValue(L.Metadata(4))
	}

	return lua.LNil
}

func (env *Environment) setExdataL(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "A":
		L.SetMetadata(0, val)
	case "B":
		L.SetMetadata(1, val)
	case "C":
		L.SetMetadata(2, val)
	case "D":
		L.SetMetadata(3, val)
	case "E":
		L.SetMetadata(4, val)
	}
}

func (env *Environment) G() *lua.LTable {
	return lua.CloneTable(env.tab._G)
}

func (env *Environment) P(fn *lua.LFunction) lua.P {
	fn.Env = env.tab._G
	cp := lua.P{
		Fn:   fn,
		NRet: lua.MultRet,
	}

	switch env.tab.mode {
	case "debug":
		cp.Protect = env.tab.protect
	default:
		cp.Protect = true
	}

	return cp
}

func (env *Environment) DoFile(L *lua.LState, path string) error {
	fn, err := L.LoadFile(path)
	if err != nil {
		return err
	}

	return L.CallByParam(env.P(fn))
}

func (env *Environment) DoString(L *lua.LState, chunk string) error {
	fn, err := L.Load(strings.NewReader(chunk), "<string>")
	if err != nil {
		return err
	}
	return L.CallByParam(env.P(fn))
}

func (env *Environment) DoChunk(L *lua.LState, chunk []byte) error {
	fn, err := L.Load(bytes.NewReader(chunk), "<chunk>")
	if err != nil {
		return err
	}
	return L.CallByParam(env.P(fn))
}

func (env *Environment) Start(co *lua.LState, v lua.VelaEntry) vela.Start {
	return tasktree.Start(co, v)
}

func (env *Environment) Call(L *lua.LState, fn *lua.LFunction, args ...lua.LValue) error {
	if L == nil {
		L = env.Coroutine()
		defer env.Free(L)
	}

	return L.CallByParam(env.P(fn), args...)
}

func (env *Environment) thread(L *lua.LState) int {
	n := L.GetTop()
	if n < 1 {
		L.RaiseError("rock.go(fn , ...) , got null")
		return 0
	}

	fn := L.CheckFunction(1)
	args := make([]lua.LValue, n-1)
	for i := 2; i <= n; i++ {
		args[i-2] = L.Get(i)
	}

	co, _ := L.NewThread()
	cp := env.P(fn)
	env.submit(func() { co.CallByParam(cp, args...) })
	return 0
}
