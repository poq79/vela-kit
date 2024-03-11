package runtime

import (
	"github.com/shirou/gopsutil/v3/process"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/worker"
	"gopkg.in/tomb.v2"
	"os"
	"runtime"
	"time"
)

const (
	Min = uint64(200) * 1024 * 1024
	Max = uint(500) * 1024 * 1024
)

func PMemory() *process.MemoryInfoStat {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		xEnv.Errorf("find process fail %v", err)
		return nil
	}

	m, err := p.MemoryInfo()
	if err != nil {
		xEnv.Errorf("find process memory fail %v", err)
		return nil
	}

	return m
}

func poll(tom *tomb.Tomb, max uint64) {
	var info runtime.MemStats
	tk := time.NewTicker(6 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-tom.Dying():
			return

		case <-tk.C:
			runtime.ReadMemStats(&info)
			if info.Alloc > max {
				xEnv.Kill(os.Kill)
				xEnv.Errorf("memory overflow %d > %d", info.Alloc, max)
				os.Exit(-1)
			}
		}

	}
}

func setMaxMemoryL(L *lua.LState) int {
	checkVM(L)

	max := uint64(L.IsInt(1)) * 1024 * 1024
	if max <= Min {
		return 0
	}

	tom := new(tomb.Tomb)
	task := func() { poll(tom, max) }
	kill := func() { tom.Kill(nil) }
	worker := worker.New(L, "memory.poll").Env(xEnv).Task(task).Kill(kill)
	xEnv.Start(L, worker).From(worker.CodeVM()).Do()
	return 0
}
