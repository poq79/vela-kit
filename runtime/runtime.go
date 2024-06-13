package runtime

import (
	"github.com/elastic/gosigar"
	"github.com/vela-ssoc/vela-kit/audit"
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

func setAgentCpuAlarmL(L *lua.LState) int {
	Times := L.IsInt(1)
	TValue := float64(L.CheckNumber(2))
	TPercent := float64(L.CheckNumber(3))
	AlarmCache := L.IsInt(4)
	AlarmInterval := L.IsInt(5)

	if Times > 0 {
		AgentAlarmCfg.Cpu.Times = Times
	}
	AgentAlarmCfg.Cpu.TValue = TValue
	if TPercent > 0 && TPercent <= 1 {
		AgentAlarmCfg.Cpu.TPercent = TPercent
	}
	if AlarmCache > 0 {
		AgentAlarmCfg.Cpu.AlarmCache = AlarmCache
	}
	if AlarmInterval > 0 {
		AgentAlarmCfg.Cpu.AlarmInterval = AlarmInterval
	}
	audit.Errorf("setAgentCpuAlarm...%s", AgentAlarmCfg.Cpu.Tojson()).From("runtime").Log().Put()
	return 0
}

func setAgentMemAlarmL(L *lua.LState) int {
	Times := L.IsInt(1)
	TValue := float64(L.CheckNumber(2))
	TPercent := float64(L.CheckNumber(3))
	AlarmCache := L.IsInt(4)
	AlarmInterval := L.IsInt(5)

	if Times > 0 {
		AgentAlarmCfg.Mem.Times = Times
	}
	AgentAlarmCfg.Mem.TValue = TValue
	if TPercent > 0 && TPercent <= 1 {
		AgentAlarmCfg.Mem.TPercent = TPercent
	}
	if AlarmCache > 0 {
		AgentAlarmCfg.Mem.AlarmCache = AlarmCache
	}
	if AlarmInterval > 0 {
		AgentAlarmCfg.Mem.AlarmInterval = AlarmInterval
	}
	audit.Errorf("setAgentMemAlarm...%s", AgentAlarmCfg.Mem.Tojson()).From("runtime").Log().Put()
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
