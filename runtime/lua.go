package runtime

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"gopkg.in/tomb.v2"
	"runtime"
)

var xEnv vela.Environment

func Constructor(env vela.Environment, callback func(v interface{}) error) {
	xEnv = env
	m := &monitor{
		tomb:     new(tomb.Tomb),
		Memory:   Min,
		Cpu:      80,
		AgentMem: uint(Min),
		AgentCpu: uint(15),
	}
	m.define(xEnv.R())
	go m.task()
	callback(m)

	rtv := lua.NewUserKV()
	rtv.Set("code", lua.NewFunction(codeL))
	rtv.Set("free", lua.NewFunction(freeL))

	rtv.Set("max_memory", lua.NewFunction(setMaxMemoryL))
	rtv.Set("max_thread", lua.NewFunction(setMaxThreadL))
	rtv.Set("agent_cpu_alarm", lua.NewFunction(setAgentCpuAlarmL))
	rtv.Set("agent_mem_alarm", lua.NewFunction(setAgentMemAlarmL))

	rtv.Set("max_cpu", lua.NewFunction(setMaxCpuL))
	rtv.Set("memory", lua.NewFunction(memoryL))
	rtv.Set("p_memory", lua.NewFunction(pMemoryL))

	rtv.Set("pprof", lua.NewFunction(pprofL))
	rtv.Set("OS", lua.S2L(runtime.GOOS))
	rtv.Set("ARCH", lua.S2L(runtime.GOARCH))
	rtv.Set("CPU_CORE", lua.LInt(runtime.NumCPU()))
	rtv.Set("windows", lua.LBool(goos == "windows"))
	rtv.Set("linux", lua.LBool(goos == "linux"))

	rtv.Set("monitor", lua.NewFunction(m.setL))

	rtv.Set("max_open_file", lua.NewFunction(SetMaxOpenFileL))

	env.Global("runtime", rtv)

}
