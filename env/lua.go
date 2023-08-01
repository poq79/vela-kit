package env

import "github.com/vela-ssoc/vela-kit/lua"

type EnvL struct {
	mode    string      //Lua 运行模式
	name    string      //lua 环境运行名称
	protect bool        //protect
	ukv     lua.UserKV  //lua 用户缓存键值对
	_G      *lua.LTable //lua 运行时字典
}

func newEnvL(mode, name string) *EnvL {
	L := lua.NewState()
	defer L.Close()

	env := &EnvL{
		mode: mode,
		name: name,
		ukv:  lua.NewUserKV(),
		_G:   L.Env,
	}
	env._G.RawSetString(name, env.ukv)
	return env
}

func (env *Environment) newEnvL(mode, name string, protect bool) {
	env.tab = newEnvL(mode, name)
	env.tab.protect = protect
}

func (env *Environment) Name() string {
	return env.tab.name
}

func (env *Environment) Set(key string, lv lua.LValue) {
	env.tab.ukv.Set(key, lv)
}

func (env *Environment) Global(key string, lv lua.LValue) {
	env.tab._G.RawSetString(key, lv)
}

func (env *Environment) Protect() bool {
	return env.tab.protect
}
