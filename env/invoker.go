package env

import (
	"context"
	"github.com/vela-ssoc/vela-kit/vela"
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
