package runtime

import (
	"errors"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// func GetCPUpct() (float64, error) {
//
//		return cpuUsage, nil
//	}
var latestCpuPct float64
var latestProcessCpuPct float64
var historyStat Stat
var historyLock sync.Mutex
var clkTck float64 = 100 // 默认100

type Stat struct {
	utime  float64
	stime  float64
	cutime float64
	cstime float64
	start  float64
	rss    float64
	uptime float64
}

func GetCurrentProcessCPUpct() (float64, float64, error) {
	pct, err := cpu.Percent(5000*time.Millisecond, false)
	if err != nil {
		return -1, -1, err
	}
	currentProcessCpuUsage, err := GetStat(os.Getpid())
	if err != nil {
		return -1, -1, err
	}
	// currentProcessCpuUsage, err := agt.CPUPercent()
	// if err != nil {
	// 	return -1, -1, err
	// }
	latestCpuPct = pct[0]
	latestProcessCpuPct = currentProcessCpuUsage

	return pct[0], currentProcessCpuUsage, nil
}

func GetCurrentProcessCPUpctLatest() (float64, float64, error) {
	// 从环形缓存中读取
	//var CPUpct float64
	//var CurrentProcessCPUpct float64
	//
	//CPUpct = 0
	//CurrentProcessCPUpct = 0
	//if rb.os.cpu.Value != nil {
	//	CPUpct = rb.os.cpu.Value.(float64)
	//}
	//if rb.agent.cpu.Value != nil {
	//	CurrentProcessCPUpct = rb.agent.cpu.Value.(float64)
	//}
	//
	//return CPUpct, CurrentProcessCPUpct, nil

	return latestCpuPct, latestProcessCpuPct, nil
}

func init() {
	historyStat = Stat{}
}

func parseFloat(val string) float64 {
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

func GetStat(pid int) (float64, error) {
	uptimeFileBytes, err := os.ReadFile(path.Join("/proc", "uptime"))
	if err != nil {
		return -1, err
	}
	uptime := parseFloat(strings.Split(string(uptimeFileBytes), " ")[0])

	procStatFileBytes, err := os.ReadFile(path.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return -1, err
	}
	splitAfter := strings.SplitAfter(string(procStatFileBytes), ")")

	if len(splitAfter) == 0 || len(splitAfter) == 1 {
		return -1, errors.New("Can't find process with this PID: " + strconv.Itoa(pid))
	}
	infos := strings.Split(splitAfter[1], " ")
	stat := &Stat{
		utime:  parseFloat(infos[12]),
		stime:  parseFloat(infos[13]),
		cutime: parseFloat(infos[14]),
		cstime: parseFloat(infos[15]),
		start:  parseFloat(infos[20]) / clkTck,
		rss:    parseFloat(infos[22]),
		uptime: uptime,
	}

	_stime := 0.0
	_utime := 0.0

	historyLock.Lock()
	defer historyLock.Unlock()

	_history := historyStat

	if _history.stime != 0 {
		_stime = _history.stime
	}

	if _history.utime != 0 {
		_utime = _history.utime
	}
	total := stat.stime - _stime + stat.utime - _utime
	total = total / clkTck

	seconds := stat.start - uptime
	if _history.uptime != 0 {
		seconds = uptime - _history.uptime
	}

	seconds = math.Abs(seconds)
	if seconds == 0 {
		seconds = 1
	}

	historyStat = *stat
	cpu := (total / seconds) * 100
	return cpu, nil
}
