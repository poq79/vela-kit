package env

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/bucket"
	"github.com/vela-ssoc/vela-kit/httpx"
	"github.com/vela-ssoc/vela-kit/logger"
	"github.com/vela-ssoc/vela-kit/require"
	"github.com/vela-ssoc/vela-kit/rtable"
	"github.com/vela-ssoc/vela-kit/runtime"
	"github.com/vela-ssoc/vela-kit/third"
	"github.com/vela-ssoc/vela-kit/variable"
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/webdav"
	"github.com/vela-ssoc/vela-kit/webutil"
	xlink "github.com/vela-ssoc/vela-xlink"
)

func (env *Environment) invoke() {

	rtable.Constructor(env, func(rt *rtable.TnlRouter) error {
		if env.router != nil {
			return fmt.Errorf("tunnel router already ok")
		}

		env.router = rt
		return nil
	})

	logger.Constructor(env, func(log vela.Log) error {
		if env.log != nil {
			return fmt.Errorf("env logger already ok")
		}

		env.log = log
		return nil
	})

	variable.Constructor(env, func(hub *variable.Hub) error {
		if env.vhu != nil {
			return fmt.Errorf("env variable hub already ok")
		}
		env.vhu = hub
		return nil
	})

	bucket.Constructor(env, func(v interface{}) error {
		if v == nil {
			return fmt.Errorf("ssc database got nil")
		}

		db, ok := v.(database)
		if !ok {
			return fmt.Errorf("invalid ssc database object , got %t", v)
		}

		env.db = db
		return nil
	})

	third.Constructor(env, func(t third.VelaThird) error {
		if env.third != nil {
			return fmt.Errorf("third object already ok")
		}
		env.third = t
		return nil
	})

	require.Constructor(env, func(p require.Pool) error {
		if env.requireHub != nil {
			return fmt.Errorf("require cache pool already ok")
		}

		env.requireHub = p
		return nil
	})

	audit.Constructor(env, func(adt *audit.Audit) error {
		if env.adt != nil {
			return fmt.Errorf("audit already ok")
		}

		env.adt = adt
		return adt.Start()
	})

	runtime.Constructor(env, func(v interface{}) error {
		if v == nil {
			return fmt.Errorf("invalid runtime object")
		}

		if m, ok := v.(monitor); ok {
			env.rtm = m
			return nil
		}

		return fmt.Errorf("invalid monitor object")
	})

	xlink.Constructor(env, func(v interface{}) {
		if v == nil {
			return
		}
		env.link = v
	})

	webdav.Constructor(env)
	webutil.Constructor(env)
	httpx.WithEnv(env)
}
