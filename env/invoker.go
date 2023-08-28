package env

import (
	"context"
	"fmt"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/bucket"
	"github.com/vela-ssoc/vela-kit/logger"
	"github.com/vela-ssoc/vela-kit/rtable"
	"github.com/vela-ssoc/vela-kit/runtime"
	"github.com/vela-ssoc/vela-kit/third"
	"github.com/vela-ssoc/vela-kit/variable"
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/webdav"
)

func (env *Environment) Logger() vela.Log {
	return env.log
}

func (env *Environment) Adt() interface{} {
	return env.adt
}

func (env *Environment) Context() context.Context {
	return env.ctx
}

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

	webdav.Constructor(env)
}
