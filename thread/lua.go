package thread

import "github.com/vela-ssoc/vela-kit/vela"

var xEnv vela.Environment

func Constructor(env vela.Environment) {
	xEnv = env
}
