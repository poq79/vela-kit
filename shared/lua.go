package shared

import (
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

func Constructor(env vela.Environment) {
	shm := &shared{data: make(map[string]*ShareBucket, 32)}

	xEnv = env
	xEnv.InitSharedEnv(shm)
	xEnv.Set("shared", shm)
}
