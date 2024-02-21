package runtime

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

// func GetCPUpct() (float64, error) {
//
//		return cpuUsage, nil
//	}
var latestCpuPct float64
var latestProcessCpuPct float64
var historyProcessTimes ProcessTimes
var historySystemTimes SystemTimes
var clkTck float64 = 100 // 默认100

type ProcessTimes struct {
	Utime  float64
	Stime  float64
	Cutime float64
	Cstime float64
	Start  float64
	Rss    float64
	Uptime float64
}

type SystemTimes struct {
	Utime   float64
	Stime   float64
	Nice    float64
	Iowait  float64
	Idle    float64
	IRQ     float64
	SoftIRQ float64
}

func GetCurrentProcessCPUpct() (float64, float64, error) {
	CpuUsage, currentProcessCpuUsage, err := GetCpuStat()
	if err != nil {
		return -1, -1, err
	}
	latestCpuPct = CpuUsage
	latestProcessCpuPct = currentProcessCpuUsage
	return CpuUsage, currentProcessCpuUsage, nil
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
	historyProcessTimes = ProcessTimes{}
}

func parseFloat(val string) float64 {
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

func getUptime() float64 {
	uptimeFileBytes, err := os.ReadFile(path.Join("/proc", "uptime"))
	if err != nil {
		return -1
	}
	uptime := parseFloat(strings.Split(string(uptimeFileBytes), " ")[0])
	return uptime
}

func (s *SystemTimes) alltime() float64 {
	return s.Stime + s.Utime + s.Nice + s.Iowait + s.IRQ + s.SoftIRQ + s.Idle
}

func GetCpuStat() (float64, float64, error) {
	uptime := getUptime()
	pid := os.Getpid()
	systemTimes, err := getSystemTimes()
	if err != nil {
		return -1, -1, err
	}
	procStatFileBytes, err := os.ReadFile(path.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return -1, -1, err
	}
	splitAfter := strings.SplitAfter(string(procStatFileBytes), ")")

	if len(splitAfter) == 0 || len(splitAfter) == 1 {
		return -1, -1, errors.New("[GetCpuStat] Can't find process with this PID: " + strconv.Itoa(pid))
	}
	infos := strings.Split(splitAfter[1], " ")
	processTimes := &ProcessTimes{
		Utime:  parseFloat(infos[12]),
		Stime:  parseFloat(infos[13]),
		Cutime: parseFloat(infos[14]),
		Cstime: parseFloat(infos[15]),
		Start:  parseFloat(infos[20]) / clkTck,
		Rss:    parseFloat(infos[22]),
		Uptime: uptime,
	}

	if historyProcessTimes.Stime == 0 && historySystemTimes.Stime == 0 {
		historyProcessTimes = *processTimes
		historySystemTimes = *systemTimes
		return 0, 0, nil
	}
	// totalDelta时间段内的进程的cpu busy时间
	processBusyDelta := processTimes.Stime + processTimes.Utime + processTimes.Cutime + processTimes.Cstime - (historyProcessTimes.Stime + historyProcessTimes.Utime + historyProcessTimes.Cutime + historyProcessTimes.Cstime)

	totalDelta := systemTimes.alltime() - historySystemTimes.alltime()
	// totalDelta时间段内的系统的cpu busy时间
	busyDelta := totalDelta - (systemTimes.Idle - historySystemTimes.Idle)

	historyProcessTimes = *processTimes
	historySystemTimes = *systemTimes
	cpu := (busyDelta / totalDelta) * 100
	processCpu := (processBusyDelta / totalDelta) * 100
	return cpu, processCpu, nil
}

func getSystemTimes() (*SystemTimes, error) {
	statFilePath := "/proc/stat"
	statContent, err := os.ReadFile(statFilePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(statContent), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "cpu" {
			times := &SystemTimes{
				Utime:   parseFloat(fields[1]),
				Nice:    parseFloat(fields[2]),
				Stime:   parseFloat(fields[3]),
				Idle:    parseFloat(fields[4]),
				Iowait:  parseFloat(fields[5]),
				IRQ:     parseFloat(fields[6]),
				SoftIRQ: parseFloat(fields[7]),
			}
			return times, nil
		}
	}

	return nil, fmt.Errorf("unable to get system CPU times")
}
