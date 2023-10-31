package denoise

import "github.com/vela-ssoc/vela-kit/vela"

var xEnv vela.Environment

func With(env vela.Environment) {
	xEnv = env
}
