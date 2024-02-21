package env

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/stdutil"
	"github.com/vela-ssoc/vela-kit/vela"
	"os"
	"os/signal"
	"syscall"
)

type closer interface {
	Name() string
	Close() error
}

func (env *Environment) start(v lua.VelaEntry) error {
	switch v.State() {

	case lua.VTRun:
		obj, ok := v.(interface{ Reload() error })
		if ok {
			return obj.Reload()
		}

		if e := v.Close(); e != nil {
			return fmt.Errorf("%s close error %v", v.Name(), e)
		}

		return v.Start()
	default:
		return v.Start()
	}
}

func (env *Environment) Register(cc vela.Closer) {
	if cc == nil {
		return
	}
	env.mbc = append(env.mbc, cc)
}

func (env *Environment) Kill(s os.Signal) {
	daemon := stdutil.New(stdutil.Daemon())
	defer func() {
		_ = daemon.Close()
	}()

	n := len(env.mbc)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		c := env.mbc[i]
		if e := c.Close(); e != nil {
			daemon.ERR("%s exit fail %v", c.Name(), e)
		} else {
			daemon.ERR("exit %v succeed signal %v", c.Name(), s)
		}
	}
}

func (env *Environment) Notify() {
	sc := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL}
	chn := make(chan os.Signal, 1)
	signal.Notify(chn, sc...)
	s := <-chn
	env.Kill(s)
}

func (env *Environment) notifyL(L *lua.LState) int {
	env.Notify()
	return 0
}

func (env *Environment) disableL(L *lua.LState) int {
	return 0
}
