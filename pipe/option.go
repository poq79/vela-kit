package pipe

import (
	"github.com/vela-ssoc/vela-kit/vela"
)

func Seek(n int) func(*Chains) {
	return func(px *Chains) {
		if n < 0 {
			return
		}
		px.seek = n
	}
}

func Env(env vela.Environment) func(*Chains) {
	return func(px *Chains) {
		px.xEnv = env
	}
}
