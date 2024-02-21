package runtime

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"syscall"
)

func SetMaxOpenFileL(L *lua.LState) int {

	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	n := L.CheckInt(1)
	if n < 1024 {
		return 0
	}

	rlimit.Max = uint64(n)
	rlimit.Cur = uint64(n)

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		L.RaiseError("%v", err)
	}

	return 0
}
