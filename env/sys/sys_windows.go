package sys

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
	"syscall"
)

func WithEnv(env vela.Environment) {
	env.Set("SIGABRT", lua.LNumber(syscall.SIGABRT))
	env.Set("SIGALRM", lua.LNumber(syscall.SIGALRM))
	env.Set("SIGBUS", lua.LNumber(syscall.SIGBUS))
	env.Set("SIGFPE", lua.LNumber(syscall.SIGFPE))
	env.Set("SIGHUP", lua.LNumber(syscall.SIGHUP))
	env.Set("SIGILL", lua.LNumber(syscall.SIGILL))
	env.Set("SIGINT", lua.LNumber(syscall.SIGINT))
	env.Set("SIGKILL", lua.LNumber(syscall.SIGKILL))
	env.Set("SIGPIPE", lua.LNumber(syscall.SIGPIPE))
	env.Set("SIGQUIT", lua.LNumber(syscall.SIGQUIT))
	env.Set("SIGSEGV", lua.LNumber(syscall.SIGSEGV))
	env.Set("SIGTERM", lua.LNumber(syscall.SIGTERM))
	env.Set("SIGTRAP", lua.LNumber(syscall.SIGTRAP))
}
