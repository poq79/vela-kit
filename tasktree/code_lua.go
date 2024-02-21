package tasktree

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/reflectx"
	"github.com/vela-ssoc/vela-kit/stdutil"
	"github.com/vela-ssoc/vela-kit/vela"
	"runtime"
)

func (cd *Code) String() string                         { return lua.B2S(cd.chunk) }
func (cd *Code) Type() lua.LValueType                   { return lua.LTObject }
func (cd *Code) AssertFloat64() (float64, bool)         { return 0, false }
func (cd *Code) AssertString() (string, bool)           { return "", false }
func (cd *Code) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (cd *Code) Peek() lua.LValue                       { return cd }

func startL(L *lua.LState) int {
	n := L.GetTop()
	for i := 1; i <= n; i++ {
		ud := L.CheckVelaData(i)
		xEnv.Start(L, ud.Data).From(ud.CodeVM()).Do()
	}
	return 0
}

func privateL(L *lua.LState) int {

	cname := L.CodeVM()
	if cname == "" {
		L.RaiseError("not allow inline , must be code vm")
		return 0
	}

	n := L.GetTop()
	for i := 1; i <= n; i++ {
		ud := L.CheckVelaData(i)
		ud.Private(L)

		root.forEach(func(key string, co *lua.LState, code *Code) bool {
			if code.inLink(cname) {
				code.ToUpdate()
				xEnv.Errorf("%s code with inline , %s set update reg", cname, code.Key())
			}
			return true
		})
	}
	return 0
}

// vela.disable("linux")
// vela.disable()

func disableL(L *lua.LState) int {
	os := L.IsString(1)
	if len(os) == 0 {
		L.RaiseError("disable code")
		return 0
	}

	if os == runtime.GOOS {
		L.RaiseError("disable code")
		return 0
	}

	return 0
}

// Index CODE结果中的PROC服务
func (cd *Code) Index(L *lua.LState, key string) lua.LValue {
	ud := cd.vela(key)
	if ud == nil {
		L.RaiseError("not found %s vela", key)
		return lua.LNil
	}

	if !ud.IsPrivate() {
		return ud
	}

	if !cd.CompareVM(L) {
		L.RaiseError("%s link %s inline vela", L.CodeVM(), ud.Data.Name())
		return lua.LNil
	}

	return ud

}

func metadataL(L *lua.LState) int {
	co, ok := CheckCodeVM(L)
	if !ok {
		return 0
	}

	key := L.CheckString(1)
	L.Push(reflectx.ToLValue(co.metadata[key], L))
	return 1
}

func paramL(L *lua.LState) int {
	co, ok := CheckCodeVM(L)
	if !ok {
		return 0
	}

	key := L.CheckString(1)
	val, ok := co.Param(key)
	if ok {
		L.Push(reflectx.ToLValue(val, L))
		return 1
	}
	return 0
}

func consoleL(L *lua.LState) int {
	chunk := lua.Format(L, 0)
	if L.Console == nil {
		console := stdutil.New(stdutil.Console())
		defer console.Close()

		console.Debug(chunk)
		return 0
	}
	L.Console.Println(chunk)
	return 0
}

func codeLuaInjectApi(env vela.Environment) {
	env.Global("start", lua.NewFunction(startL))
	env.Global("private", lua.NewFunction(privateL))
	env.Global("metadata", lua.NewFunction(metadataL))
	env.Global("console", lua.NewFunction(consoleL))
	env.Set("disable", lua.NewFunction(disableL))
	env.Set("param", lua.NewFunction(paramL))
}
