package runtime

import (
	"github.com/elastic/gosigar"
	"github.com/vela-ssoc/vela-kit/lua"
	"os"
	"runtime"
	"runtime/debug"
)

const (
	code = "runtime"
	goos = runtime.GOOS
)

func codeL(L *lua.LState) int {
	L.Push(lua.S2L(L.CodeVM()))
	return 1
}

func checkVM(L *lua.LState) {
	if L.CodeVM() != code {
		L.RaiseError("not allow with %s , must %s", L.CodeVM(), code)
	}
}

func freeL(_ *lua.LState) int {
	debug.FreeOSMemory()
	return 0
}

func setMaxCpuL(L *lua.LState) int {
	checkVM(L)

	n := L.IsInt(1)
	if n >= runtime.NumCPU() || n <= 0 {
		return 0
	}

	runtime.GOMAXPROCS(n)
	return 0
}

func setMaxThreadL(L *lua.LState) int {
	checkVM(L)

	n := L.IsInt(1)
	if n <= 0 {
		return 0
	}

	debug.SetMaxThreads(n)
	return 0
}

func memoryL(L *lua.LState) int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	L.Push(lua.LNumber(m.Alloc))
	L.Push(lua.LNumber(m.TotalAlloc))
	L.Push(lua.LNumber(m.Sys))
	return 3
}

func pMemoryL(L *lua.LState) int {
	pid := os.Getpid()
	mem := gosigar.ProcMem{}
	err := mem.Get(pid)
	if err != nil {
		xEnv.Errorf("find process fail %v", err)
		L.Push(lua.LInt(-1))
		return 0
	}

	lv := lua.NewMap(4, false)
	lv.Set("pid", lua.LNumber(pid))
	lv.Set("size", lua.LNumber(mem.Size))
	lv.Set("rss", lua.LNumber(mem.Resident))
	lv.Set("share", lua.LNumber(mem.Share))
	L.Push(lv)
	return 1
}
