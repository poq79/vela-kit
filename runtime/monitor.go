package runtime

import (
	"container/ring"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"gopkg.in/tomb.v2"
	"os"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"
)

type Ring struct {
	cpu *ring.Ring
	mem *ring.Ring
}

type RingBuffer struct {
	os    *Ring
	agent *Ring
}

var rb = &RingBuffer{
	os: &Ring{
		cpu: ring.New(120),
		mem: ring.New(120),
	},

	agent: &Ring{
		cpu: ring.New(120),
		mem: ring.New(120),
	},
}

type Cache struct {
	CPU   float64
	Agent float64
	Quiet time.Time
}

type monitor struct {
	doing    uint32
	counter  uint32
	tomb     *tomb.Tomb
	Memory   uint64
	CPU      uint
	CPUTimes []cpu.TimesStat
	AgentCpu uint
	AgentMem uint
	State    Cache
	Day      int
	Action   string
}

func (m *monitor) Index(L *lua.LState, key string) lua.LValue {
	return lua.LNil
}

func (m *monitor) ttl() int64 {
	if m.Day < 0 {
		return int64(86400 * 7)
	}
	return int64(m.Day * 86400)
}

func (m *monitor) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "agent_mem":
		n := lua.IsInt(val)
		if n <= 10 {
			m.AgentMem = 10
			return
		}
		m.AgentMem = uint(n)

	case "agent_cpu":
		n := lua.IsInt(val)
		if n <= 15 {
			m.AgentCpu = 15
		}

		m.AgentCpu = uint(n)

	case "action":
		m.Action = val.String()

	case "memory":
		n := lua.IsInt(val)
		if uint64(n) <= Min {
			m.Memory = Min
			return
		}
		m.Memory = uint64(n)
	}
}

func (m *monitor) Disable() {
	now := time.Now()
	if now.Unix()-m.State.Quiet.Unix() > 3600 {
		m.State.Quiet = time.Now()
		return
	}
	m.State.Quiet = time.Now()
}

func (m *monitor) Quiet() bool {
	now := time.Now().Unix()

	if now-m.State.Quiet.Unix() < 3600 {
		return true
	}
	return false
}

func (m *monitor) Decision() {
	switch m.Action {
	case "exit":
		os.Exit(-1)
	case "quiet":
		m.Disable()
	default:
		//xEnv.QuietOn()
	}
}

func (m *monitor) AgentAlloc() {
	var info runtime.MemStats
	runtime.ReadMemStats(&info)

	if info.Alloc > m.Memory {
		audit.Errorf("memory overflow %d > %d", info.Alloc, m.Memory).From("runtime").Log().Put()
		debug.FreeOSMemory()
	}
}

func (m *monitor) store(key string, v float64) {
	r := &Record{
		ID:    time.Now().Unix(),
		Value: v,
	}

	switch key {
	case "os.cpu":
		m.State.CPU = v
		rb.os.cpu.Value = r
		rb.os.cpu = rb.os.cpu.Next()
	case "os.mem":
		rb.os.mem.Value = r
		rb.os.mem = rb.os.mem.Next()
	case "agent.cpu":
		m.State.Agent = v //agent cpu
		rb.agent.cpu.Value = r
		rb.agent.cpu = rb.agent.cpu.Next()
	case "agent.mem":
		rb.agent.mem.Value = r
		rb.agent.mem = rb.agent.mem.Next()
	}
}

// 弃用, 系统时间计算和Agent时间一起算, 在agent()函数中实现
//func (m *monitor) Cpu() {
//	t2, err := cpu.Times(false)
//	if err != nil {
//		return
//	}
//
//	t1 := m.CPUTimes
//	m.CPUTimes = t2
//
//	if len(t1) == 0 {
//		return
//	}
//
//	ret, err := CalculateAllBusy(t1, t2)
//	if err != nil {
//		xEnv.Errorf("calculate cpu fail %v", err)
//		return
//	}
//
//	//xEnv.Errorf("cpu avg:%v    max:%v    usage:%v", avg(ret), max(ret), ret)
//	m.store("os.cpu", ret[0])
//}

func (m *monitor) Mem() {
	mi, err := mem.VirtualMemory()
	if err != nil {
		xEnv.Errorf("runtime mem monitor %v", err)
		return
	}

	//xEnv.Errorf("cpu avg:%v    max:%v    usage:%v", avg(ret), max(ret), ret)
	m.store("os.mem", mi.UsedPercent)
}

func (m *monitor) Agent() {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		xEnv.Errorf("ps agent fail %v", err)
		return
	}

	//cpct, err := p.CPUPercent()
	//if err != nil {
	//	xEnv.Errorf("ps agent cpu percent fail %v", err)
	//}
	//m.store("agent.cpu", cpct)
	//cpct, err := AgentCpu(int32(os.Getpid()))
	//systemCpuUsage, currentProcessCpuUsage, err := GetCurrentProcessCPUpct(nil)

	systemCpuUsage, currentProcessCpuUsage, err := GetCurrentProcessCPUpct()

	if err != nil {
		xEnv.Errorf("ps agent cpu percent fail %v", err)
		return
	}

	m.store("os.cpu", systemCpuUsage)
	m.store("agent.cpu", currentProcessCpuUsage)
	mpct, err := p.MemoryPercent()
	if err != nil {
		xEnv.Errorf("ps agent mem percent fail %v", err)
	}
	m.store("agent.mem", float64(mpct))
}

func (m *monitor) inc() {
	atomic.AddUint32(&m.counter, 1)
}

func (m *monitor) Chance(r *ring.Ring, e float64, h float64) bool { // 120次 , 90次 , 0.8

	var total int
	var exceed int
	r.Do(func(a interface{}) {
		if a == nil {
			return
		}

		entry, ok := a.(*Record)
		if !ok {
			return
		}
		total++

		if entry.Value >= e {
			exceed++
		}
	})

	if total < 120 {
		return false
	}

	if float64(exceed)/float64(total) > h {
		return true
	}

	return false
}

func (m *monitor) exec() bool {

	if m.Chance(rb.agent.cpu, 20, 0.7) {
		audit.Errorf("agent.cpu overflow  2min").From("runtime").Log().Put()
		m.Decision()
	}

	if m.Chance(rb.agent.mem, float64(Min), 0.8) {
		audit.Errorf("agent.mem overflow 150M 2min").From("runtime").Log().Put()
		m.Decision()
	}

	return false
}

func (m *monitor) task() {
	tk := time.NewTicker(5 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-m.tomb.Dying():
			return

		case <-tk.C:
			m.inc()
			m.Agent()
			m.AgentAlloc()
			//m.Cpu()
			m.Mem()
			m.exec()
		}
	}

}

func (m *monitor) kill() {
	m.tomb.Kill(fmt.Errorf("runtime over"))
}

func (m *monitor) setL(L *lua.LState) int {
	checkVM(L)
	cfg := L.CheckTable(1)
	cfg.Range(func(key string, val lua.LValue) {
		m.NewIndex(L, key, val)
	})
	return 0
}
