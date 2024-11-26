package runtime

import (
	"container/ring"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	risk "github.com/vela-ssoc/vela-risk"
	vtag "github.com/vela-ssoc/vela-tag"
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
	Cpu      uint
	CPUTimes []cpu.TimesStat
	AgentCpu uint
	AgentMem uint
	State    Cache
	Day      int
	Action   string
}

type agentAlarmCfg struct {
	Cpu alarmCfg `json:"Cpu"`
	Mem alarmCfg `json:"Mem"`
}
type alarmCfg struct {
	// Times tValue tPercent
	Times         int     `json:"Times"`         //往前取多少次数据
	TValue        float64 `json:"TValue"`        //告警阈值
	TPercent      float64 `json:"TPercent"`      //达到阈值的百分比
	AlarmCache    int     `json:"AlarmCache"`    //告警计数器缓存的时间(秒)
	AlarmInterval int     `json:"AlarmInterval"` //告警间隔 如10次则是连续10次告警才上传告警一次
}

// 往前追溯Times个监控周期内 有TPercent的agentcpu占用数据大于TValue%,则产生告警
// 但是在AlarmCache秒内每产生AlarmInterval次才上传告警一次
func (a alarmCfg) Tojson() string {
	return fmt.Sprintf(`{"Times":%d,"TValue":%f,"TPercent":%f,"AlarmCache":%d,"Alarm Interval":%d}`, a.Times, a.TValue, a.TPercent, a.AlarmCache, a.AlarmInterval)
}

var AgentAlarmCfg = agentAlarmCfg{
	Cpu: alarmCfg{
		Times:         24,
		TValue:        20,
		TPercent:      0.9,
		AlarmCache:    3600,
		AlarmInterval: 120,
	},
	Mem: alarmCfg{
		Times:         24,
		TValue:        float64(Min),
		TPercent:      0.9,
		AlarmCache:    3600,
		AlarmInterval: 120,
	},
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
		//audit.Errorf("memory overflow %d > %d", info.Alloc, m.Memory).From("runtime").Log().Put()
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
//
//	func (m *monitor) Cpu() {
//		t2, err := cpu.Times(false)
//		if err != nil {
//			return
//		}
//
//		t1 := m.CPUTimes
//		m.CPUTimes = t2
//
//		if len(t1) == 0 {
//			return
//		}
//
//		ret, err := CalculateAllBusy(t1, t2)
//		if err != nil {
//			xEnv.Errorf("calculate cpu fail %v", err)
//			return
//		}
//
//		//xEnv.Errorf("cpu avg:%v    max:%v    usage:%v", avg(ret), max(ret), ret)
//		m.store("os.cpu", ret[0])
//	}
func (m *monitor) CPU() float64 {
	return latestCpuPct
}

func (m *monitor) AgentCPU() float64 {
	return latestProcessCpuPct
}

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

// Chance 监控告警触发器  最近的times次采集数据中有tPercent的比例的数据大于tValue, 返回true
func (m *monitor) Chance(r *ring.Ring, times int, tValue float64, tPercent float64) bool { // 120次 , 90次 , 0.8

	var total int
	var exceed int

	// 往前回溯times次数据
	for p := r.Prev(); p != r && total < times; p = p.Prev() {
		if p.Value == nil {
			break
		}
		entry, ok := p.Value.(*Record)
		if !ok {
			break
		}
		total++

		if entry.Value >= tValue {
			exceed++
		}
	}

	if total < times {
		return false
	}
	if float64(exceed)/float64(total) > tPercent {
		return true
	}

	return false
}

// agent自身资源占用监控告警
func (m *monitor) exec() bool {
	shmCpu := xEnv.Shm("SHM-RUNTIME-AGENT-CPU", "OVERFLOW")
	shmMem := xEnv.Shm("SHM-RUNTIME-AGENT-MEM", "OVERFLOW")
	t := vtag.NewTag()

	// 默认 24*5秒=120秒
	if m.Chance(rb.agent.cpu, AgentAlarmCfg.Cpu.Times, AgentAlarmCfg.Cpu.TValue, AgentAlarmCfg.Cpu.TPercent) {
		if n, _ := shmCpu.Incr("acpu.alert", 1, AgentAlarmCfg.Cpu.AlarmCache*1000); n > 1 {
			if n%AgentAlarmCfg.Cpu.AlarmInterval == 0 {
				audit.Errorf("agent.cpu overflow %ds (last %ds alert times:%d)", AgentAlarmCfg.Cpu.Times*5, AgentAlarmCfg.Cpu.AlarmCache, n).From("runtime").Log().Put()
				rev := risk.Monitor()
				rev.Metadata = make(map[string]lua.LValue)
				rev.FromCode = "runtime"
				rev.Metadata["config"] = lua.S2L(AgentAlarmCfg.Cpu.Tojson())
				rev.Subject = "agent.cpu overflow"
				rev.Payload = fmt.Sprintf("agent.cpu overflow %ds (last %ds alert times:%d)", AgentAlarmCfg.Cpu.Times*5, AgentAlarmCfg.Cpu.AlarmCache, n)
				rev.Send()
			}
		} else {
			audit.Errorf("agent.cpu overflow %ds", AgentAlarmCfg.Mem.Times*5).From("runtime").Log().Alert().Put()
			t.AddTag("acpu_overflow")
			t.Send()
		}
		m.Decision()
	}
	if m.Chance(rb.agent.mem, AgentAlarmCfg.Mem.Times, AgentAlarmCfg.Mem.TValue, AgentAlarmCfg.Mem.TPercent) {
		if n, _ := shmMem.Incr("amem.alert", 1, AgentAlarmCfg.Mem.AlarmCache*1000); n > 1 {
			if n%AgentAlarmCfg.Mem.AlarmInterval == 0 {
				audit.Errorf("agent.mem overflow %ds (last %ss alert times:%d)", AgentAlarmCfg.Mem.Times*5, AgentAlarmCfg.Mem.AlarmCache, n).From("runtime").Log().Put()
				rev := risk.Monitor()
				rev.Metadata = make(map[string]lua.LValue)
				rev.FromCode = "runtime"
				rev.Metadata["config"] = lua.S2L(AgentAlarmCfg.Cpu.Tojson())
				rev.Subject = "agent.mem overflow"
				rev.Payload = fmt.Sprintf("agent.mem overflow %ds (last %ss alert times:%d)", AgentAlarmCfg.Mem.Times*5, AgentAlarmCfg.Mem.AlarmCache, n)
				rev.Send()
			}
		} else {
			audit.Errorf("agent.mem overflow %ds", AgentAlarmCfg.Mem.Times*5).From("runtime").Log().Alert().Put()
			t.AddTag("amem_overflow")
			t.Send()
		}
		m.Decision()
	}

	return false
}

func (m *monitor) uptime() time.Duration {
	return time.Duration(m.counter) * time.Second * 5
}

func (m *monitor) clean() {
	fact := 86400 * time.Second

	if m.counter <= 1 {
		goto Handle

	}

	if m.uptime()%fact != 0 {
		return
	}

Handle:
	cli := xEnv.R().Cli()
	_, _ = cli.Get("/clean/agent/logger?space=1073741824")    //1G
	_, _ = cli.Get("/clean/audit/logger?space=524288000")     //500M
	_, _ = cli.Get("/clean/agent/executable?space=524288000") //50M
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
			m.clean()
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
