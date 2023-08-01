package plugin

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"os"
	"plugin"
	"runtime"
)

var (
	xEnv vela.Environment
)

func check(L *lua.LState) string {
	return L.CheckString(1)
}

func open(L *lua.LState, path string) int {
	//判断文件是否存在
	now, err := os.Stat(path)
	if err != nil {
		L.RaiseError("%s read stat fail %v", path, err)
		return 0
	}

	//获取缓存
	key := fmt.Sprintf("linux_plugin_%s", path)
	lv, ok := xEnv.Find(key)
	if !ok {
		return load(L, key, path, now)
	}

	//解析对象 并判断文件是否变动
	old := lv.(os.FileInfo)
	if old.ModTime() == now.ModTime() {
		xEnv.Infof("%s plugin running", path)
		return 0
	}

	audit.Errorf("restart agent with %s plugin old=%d now=%d",
		path, old.ModTime().Unix(), now.ModTime().Unix())

	os.Exit(-1)
	return 0
}

func load(L *lua.LState, key string, path string, stat os.FileInfo) int {
	p, err := plugin.Open(path)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	sym, err := p.Lookup("WithEnv")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	sym.(func(vela.Environment))(xEnv)
	xEnv.Store(key, stat)

	return 0

}

func newLuaPlugin(L *lua.LState) int {
	path := check(L)
	return open(L, path)
}

func Constructor(env vela.Environment) {
	xEnv = env
	xEnv.Infof("plugin running in %s", runtime.GOOS)
	env.Set("plugin", lua.NewFunction(newLuaPlugin))
}
