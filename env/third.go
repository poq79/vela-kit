package env

import (
	"github.com/vela-ssoc/vela-kit/vela"
)

func (env *Environment) ThirdInfo(name string) *vela.ThirdInfo {
	return env.third.Info(name)
}

func (env *Environment) Third(name string) (*vela.ThirdInfo, error) {
	return env.third.Load(name)
}
