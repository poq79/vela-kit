package runtime

import (
	"fmt"
	"math"
	"syscall"
	"unsafe"
)

var (
	modkernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGetSystemTimes  = modkernel32.NewProc("GetSystemTimes")
	procGetProcessTimes = modkernel32.NewProc("GetProcessTimes")
)

var (
	latestCpuPct        float64
	latestProcessCpuPct float64
	m_agentCpuTime      int64 = 0

	m_preidleTime   = FileTime{}
	m_prekernelTime = FileTime{}
	m_preuserTime   = FileTime{}
)

type FileTime struct {
	dwLowDateTime  uint64
	dwHighDateTime uint64
}

func CompareFileTime2(time1 FileTime, time2 FileTime) int64 {
	a := int64((time1.dwHighDateTime)<<32 | time1.dwLowDateTime)
	b := int64((time2.dwHighDateTime)<<32 | time2.dwLowDateTime)
	return b - a
}

func GetCPUpct() (float64, error) {
	var idleTime, kernelTime, userTime FileTime

	// 调用GetSystemTimes函数
	ret, _, err := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)
	if ret == 0 {
		fmt.Println("Error:", err)
		return -1, err
	}

	idle := CompareFileTime2(m_preidleTime, idleTime)
	kernel := CompareFileTime2(m_prekernelTime, kernelTime)
	user := CompareFileTime2(m_preuserTime, userTime)
	var cpuUsage float64
	cpuUsage = 0
	if kernel+user == 0 {
		cpuUsage = 0
	} else {
		//（总的时间-空闲时间）/总的时间=占用cpu的时间就是使用率
		cpuUsage = math.Abs(float64((kernel+user-idle)*100) / float64(kernel+user))
	}
	m_preidleTime = idleTime
	m_prekernelTime = kernelTime
	m_preuserTime = userTime

	return cpuUsage, nil
}

func GetCurrentProcessCPUpct() (float64, float64, error) {
	// 获取系统cpu时间
	var idleTime, kernelTime, userTime FileTime

	// 调用GetSystemTimes函数
	ret, _, err := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)
	if ret == 0 {
		fmt.Println("Error:", err)
		return -1, -1, err
	}

	// 获取进程句柄
	processHandle, err := syscall.GetCurrentProcess()
	if err != nil {
		return -1, -1, err
	}
	defer syscall.CloseHandle(processHandle)

	// 获取当前进程的 Cpu 时间
	var creationTimeAgent, exitTimeAgent, kernelTimeAgent, userTimeAgent FileTime

	// 获取进程时间信息
	err = syscall.GetProcessTimes(processHandle,
		(*syscall.Filetime)(unsafe.Pointer(&creationTimeAgent)),
		(*syscall.Filetime)(unsafe.Pointer(&exitTimeAgent)),
		(*syscall.Filetime)(unsafe.Pointer(&kernelTimeAgent)),
		(*syscall.Filetime)(unsafe.Pointer(&userTimeAgent)))
	if err != nil {
		fmt.Println("无法获取当前进程cpu时间:", err)
		return -1, -1, err
	}
	// 计算 Cpu 时间
	agentCpuTime := int64((kernelTimeAgent.dwHighDateTime)<<32|kernelTimeAgent.dwLowDateTime) + int64((userTimeAgent.dwHighDateTime)<<32|userTimeAgent.dwLowDateTime)

	idle := CompareFileTime2(m_preidleTime, idleTime)
	kernel := CompareFileTime2(m_prekernelTime, kernelTime)
	user := CompareFileTime2(m_preuserTime, userTime)
	agent := agentCpuTime - m_agentCpuTime
	var cpuUsage float64
	var currentProcessCpuUsage float64

	cpuUsage = 0
	currentProcessCpuUsage = 0
	if kernel+user == 0 {
		cpuUsage = 0
	} else {
		//（总的时间-空闲时间）/总的时间=占用cpu的时间就是使用率
		cpuUsage = math.Abs(float64((kernel+user-idle)*100) / float64(kernel+user))
		currentProcessCpuUsage = math.Abs(float64(agent*100) / float64(kernel+user))
	}
	m_preidleTime = idleTime
	m_prekernelTime = kernelTime
	m_preuserTime = userTime
	m_agentCpuTime = agentCpuTime

	latestCpuPct = cpuUsage
	latestProcessCpuPct = currentProcessCpuUsage
	return cpuUsage, currentProcessCpuUsage, nil
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
