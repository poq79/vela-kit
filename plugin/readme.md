# plugin
兼容golang的plugin方法

**注意这里加载的so必须实现Constructor(xcall.Env)方法**

## 使用模块
先要生成对应的so文件
go build -buildmode=plugin demo.go

```go
    package main

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/xreflect"
)

type config struct {
	Name string `lua:"name"    type:"string"`
	Age  int    `lua:"age,18"  type:"int"`
}

func (cfg *config) Push(L *lua.LState) int {
	L.CheckString(1)
	return 0
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{}

	if e := xreflect.ToStruct(tab, cfg); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return cfg
}

func newLuaDemo(L *lua.LState) int {
	cfg := newConfig(L)
	L.Push(xreflect.ToLValue(cfg, L))
	return 1
}

func WithEnv(env assert.Environment) {
	env.Set("demo", lua.NewFunction(newLuaDemo))

}
```

## 配置脚本
```lua
    local demo = vela.plugin("xxx.so") --路径
    demo.debug()
```