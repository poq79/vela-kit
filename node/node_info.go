package node

import (
	"encoding/json"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/valyala/fasthttp"
	vdisk "github.com/vela-ssoc/vela-disk"
	vruntime "github.com/vela-ssoc/vela-kit/runtime"
)

type Info struct {
	Inet         string                  `json:"inet"`
	Inet6        string                  `json:"inet6"`
	Mac          string                  `json:"mac"`
	HostId       string                  `json:"host_id"`
	Hostname     string                  `json:"hostname"`
	Release      string                  `json:"release"`
	Version      string                  `json:"version"`
	Family       string                  `json:"family"`
	Uptime       uint64                  `json:"uptime"`
	BootAt       uint64                  `json:"boot_at"`
	Virtual      string                  `json:"virtual"`
	VirtualRole  string                  `json:"virtual_role"`
	ProcNumber   int                     `json:"proc_number"`
	MemTotal     uint64                  `json:"mem_total"`
	MemFree      uint64                  `json:"mem_free"`
	MemUsed      uint64                  `json:"mem_used"`
	MemAvailable uint64                  `json:"mem_available"`
	MemPct       float64                 `json:"mem_pct"`
	SwapTotal    uint64                  `json:"swap_total"`
	SwapInPages  uint64                  `json:"swap_in_pages"`
	SwapOutPages uint64                  `json:"swap_out_pages"`
	SwapFree     uint64                  `json:"swap_free"`
	Arch         string                  `json:"arch"`
	CpuCore      int                     `json:"cpu_core"`
	CpuInfo      []cpu.InfoStat          `json:"cpu_info"`
	CpuPct       float64                 `json:"cpu_pct"`
	AgentTotal   uint64                  `json:"agent_total"`
	AgentAlloc   uint64                  `json:"agent_alloc"`
	AgentMem     *process.MemoryInfoStat `json:"agent_mem"`
	AgentCPU     float64                 `json:"agent_cpu"`
	Disk         []vdisk.Disk            `json:"disk_info"`
}

func (i *Info) Byte() []byte {
	chunk, err := json.Marshal(i)
	if err != nil {
		xEnv.Errorf("node info marshal fail %v", err)
		return nil
	}

	return chunk
}

func (i *Info) Mem() error {
	st, err := mem.VirtualMemory()
	if err != nil {
		xEnv.Errorf("get memory stats error: %v", err)
		return err
	}

	i.MemTotal = st.Total
	i.MemFree = st.Free
	i.MemUsed = st.Used
	i.MemPct = st.UsedPercent
	i.MemAvailable = st.Available
	return nil
}

func (i *Info) Cpu() error {
	v, err := cpu.Info()
	if err != nil {
		return err
	}
	i.CpuInfo = v
	return nil
}

func (i *Info) Agt() error {
	pid := os.Getpid()
	agt, err := process.NewProcess(int32(pid))
	if err != nil {
		return err
	}
	// cpu, err := agt.CPUPercent()
	// 实时请求的接口,采集时间不固定, 就不单独做一次采集了, 直接读历史数据(5秒内)
	// 本来就不准, 而且还干扰正常定时采集的数据
	systemCpuUsage, currentProcessCpuUsage, err := vruntime.GetCurrentProcessCPUpctLatest()
	if err != nil {
		xEnv.Errorf("agent process got cpu percent fail %v", err)
	}
	i.CpuPct = systemCpuUsage
	i.AgentCPU = currentProcessCpuUsage
	// fmt.Printf("systemCpu:%.2f%%  Agt PID %d,cpu:%.2f%%\n", systemCpuUsage, pid, currentProcessCpuUsage)
	// fmt.Println("Agt PID", pid)
	mem, err := agt.MemoryInfo()
	if err != nil {
		xEnv.Errorf("agent process got mem percent fail %v", err)
		i.AgentMem = &process.MemoryInfoStat{}
	} else {
		i.AgentMem = mem
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	i.AgentAlloc = m.Alloc
	i.AgentTotal = m.TotalAlloc

	return nil
}

func (i *Info) DiskInfo(diskSum *vdisk.Summary) error {
	if diskSum == nil {
		diskSum = vdisk.New()
	}

	diskSum.Update()
	for n := 0; n < len(diskSum.Device); n++ {
		if vdisk.IsIgnorePartition(diskSum.Device[n]) {
			continue
		}
		i.Disk = append(i.Disk, *diskSum.Device[n])
	}
	return nil
}

func (i *Info) Swap() error {
	swap, err := mem.SwapMemory()
	if err != nil {
		xEnv.Errorf("got swap memory fail %v", err)
		return err
	}

	i.SwapTotal = swap.Total
	i.SwapInPages = swap.PgIn
	i.SwapOutPages = swap.PgOut
	i.SwapFree = swap.Free
	return nil
}

func (nd *node) Info(ctx *fasthttp.RequestCtx) error {
	ident := xEnv.Ident()
	hi, err := host.Info()
	diskSum := vdisk.New()
	if err != nil {
		return err
	}

	v := &Info{
		Inet:        ident.Inet.String(),
		Mac:         ident.MAC,
		HostId:      hi.HostID,
		Hostname:    hi.Hostname,
		Release:     hi.Platform,
		Version:     hi.KernelVersion,
		Uptime:      hi.Uptime,
		BootAt:      hi.BootTime,
		Virtual:     hi.VirtualizationSystem,
		VirtualRole: hi.VirtualizationRole,
		ProcNumber:  int(hi.Procs),
		Arch:        runtime.GOARCH,
		CpuCore:     runtime.NumCPU(),
		CpuPct:      xEnv.CPU(),
		AgentCPU:    xEnv.AgentCPU(),
	}

	if e := v.Mem(); e != nil {
		xEnv.Errorf("node memory info got fail %v", e)
	}

	if e := v.Swap(); e != nil {
		xEnv.Errorf("node swap info got fail %v", e)
	}

	// Agt()里面集成了系统cpu资源占用监控, 这里获取的是cpu硬件基础信息
	if e := v.Cpu(); e != nil {
		xEnv.Errorf("node cpu info got fail %v", e)
	}

	if e := v.Agt(); e != nil {
		xEnv.Errorf("node agent info got fail %v", e)
	}

	if e := v.DiskInfo(diskSum); e != nil {
		xEnv.Errorf("node disk info got fail %v", e)
	}
	ctx.Write(v.Byte())
	return nil
}
