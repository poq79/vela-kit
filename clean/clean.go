package clean

import "github.com/vela-ssoc/vela-kit/vela"

var xEnv vela.Environment

// 删除各种日志

func WithEnv(env vela.Environment) {
	xEnv = env
}
