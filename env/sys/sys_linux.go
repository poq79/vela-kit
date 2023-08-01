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
	env.Set("SIGCHLD", lua.LNumber(syscall.SIGCHLD))
	env.Set("SIGCLD", lua.LNumber(syscall.SIGCLD))
	env.Set("SIGCONT", lua.LNumber(syscall.SIGCONT))
	env.Set("SIGFPE", lua.LNumber(syscall.SIGFPE))
	env.Set("SIGHUP", lua.LNumber(syscall.SIGHUP))
	env.Set("SIGILL", lua.LNumber(syscall.SIGILL))
	env.Set("SIGINT", lua.LNumber(syscall.SIGINT))
	env.Set("SIGIO", lua.LNumber(syscall.SIGIO))
	env.Set("SIGIOT", lua.LNumber(syscall.SIGIOT))
	env.Set("SIGKILL", lua.LNumber(syscall.SIGKILL))
	env.Set("SIGPIPE", lua.LNumber(syscall.SIGPIPE))
	env.Set("SIGPOLL", lua.LNumber(syscall.SIGPOLL))
	env.Set("SIGPROF", lua.LNumber(syscall.SIGPROF))
	env.Set("SIGPWR", lua.LNumber(syscall.SIGPWR))
	env.Set("SIGQUIT", lua.LNumber(syscall.SIGQUIT))
	env.Set("SIGSEGV", lua.LNumber(syscall.SIGSEGV))
	env.Set("SIGSTKFLT", lua.LNumber(syscall.SIGSTKFLT))
	env.Set("SIGSTOP", lua.LNumber(syscall.SIGSTOP))
	env.Set("SIGSYS", lua.LNumber(syscall.SIGSYS))
	env.Set("SIGTERM", lua.LNumber(syscall.SIGTERM))
	env.Set("SIGTRAP", lua.LNumber(syscall.SIGTRAP))
	env.Set("SIGTSTP", lua.LNumber(syscall.SIGTSTP))
	env.Set("SIGTTIN", lua.LNumber(syscall.SIGTTIN))
	env.Set("SIGTTOU", lua.LNumber(syscall.SIGTTOU))
	env.Set("SIGUNUSED", lua.LNumber(syscall.SIGUNUSED))
	env.Set("SIGURG", lua.LNumber(syscall.SIGURG))
	env.Set("SIGUSR1", lua.LNumber(syscall.SIGUSR1))
	env.Set("SIGUSR2", lua.LNumber(syscall.SIGUSR2))
	env.Set("SIGVTALRM", lua.LNumber(syscall.SIGVTALRM))
	env.Set("SIGWINCH", lua.LNumber(syscall.SIGWINCH))
	env.Set("SIGXCPU", lua.LNumber(syscall.SIGXCPU))
	env.Set("SIGXFSZ", lua.LNumber(syscall.SIGXFSZ))
}
