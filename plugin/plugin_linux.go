package plugin

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>

static uintptr_t pluginOpen(const char* path, char** err) {
	void* h = dlopen(path, RTLD_NOW|RTLD_GLOBAL);
	if (h == NULL) {
		*err = (char*)dlerror();
	}
	return (uintptr_t)h;
}

static void* pluginLookup(uintptr_t h, const char* name, char** err) {
	void* r = dlsym((void*)h, name);
	if (r == NULL) {
		*err = (char*)dlerror();
	}
	return r;
}

typedef struct { void *t; void *v; } Interface;

static Interface Setup(void* f, Interface *v1) {
	Interface r;
	Interface (*vela_setup_fn)(Interface);
	vela_setup_fn = (Interface (*)(Interface))f;
	r = vela_setup_fn(*v1);
	return r;
}
*/
import "C"
import (
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"os"
	"plugin"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

var (
	xEnv vela.Environment
)

type Entry struct {
	State os.FileInfo
	Func  func(vela.Environment)
}

type Program struct {
	mu    sync.Mutex
	entry map[string]*Entry
}

func (p *Program) check(L *lua.LState) string {
	return L.CheckString(1)
}

func (p *Program) Exec(L *lua.LState, key string, path string, stat os.FileInfo) int {
	pl, err := plugin.Open(path)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	sym, err := pl.Lookup("WithEnv")
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	fn := sym.(func(vela.Environment))
	p.entry[key] = &Entry{
		State: stat,
		Func:  fn,
	}

	fn(xEnv)
	return 0
}

func (p *Program) open(L *lua.LState, path string) int {
	//判断文件是否存在
	now, err := os.Stat(path)
	if err != nil {
		L.RaiseError("%s read stat fail %v", path, err)
		return 0
	}

	//获取缓存
	p.mu.Lock()
	defer p.mu.Unlock()

	key := "linux_plugin_" + strings.ReplaceAll(path, "/", "_")
	lv, ok := p.entry[key]
	if !ok {
		return p.Exec(L, key, path, now)
	}

	//解析对象 并判断文件是否变动
	if lv.State.ModTime() == now.ModTime() {
		xEnv.Infof("%s plugin running", path)
		return 0
	}

	audit.Errorf("restart agent with %s plugin old=%d now=%d",
		path, lv.State.ModTime().Unix(), now.ModTime().Unix()).From(L.CodeVM()).Put()

	xEnv.Kill(os.Kill)
	os.Exit(-1)
	return 0

}

func (p *Program) Load(L *lua.LState) int {
	path := p.check(L)
	return p.open(L, path)
}

func (p *Program) syscallL(L *lua.LState) int {
	name := L.CheckString(1)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	hd := C.dlopen(cname, C.RTLD_LAZY)
	if hd == nil {
		L.RaiseError("dlopen fail")
		return 0
	}

	//入口
	sn := C.CString("Setup")
	defer C.free(unsafe.Pointer(sn))

	sym := C.dlsym(hd, sn)
	if sym == nil {
		L.RaiseError("not found setup")
		return 0
	}

	v2 := C.Setup(sym, (*C.Interface)(unsafe.Pointer(&xEnv)))
	vv := *(*any)(unsafe.Pointer(&v2))

	switch lv := vv.(type) {
	case lua.LValue:
		L.Push(lv)
	default:
		ad := lua.NewAnyData(lv, lua.Reflect(lua.OFF))
		L.Push(ad)
	}

	return 1
}

func (p *Program) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "syscall":
		return lua.NewFunction(p.syscallL)
	}
	return lua.LNil
}

func Constructor(env vela.Environment) {
	xEnv = env
	xEnv.Infof("plugin running in %s", runtime.GOOS)
	p := &Program{
		entry: make(map[string]*Entry),
	}

	env.Set("plugin", lua.NewExport("lua.plugin.export", lua.WithFunc(p.Load), lua.WithIndex(p.Index)))
}
