package env

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/vela"
)

type shared interface {
	NewLRU(string, int) vela.Shared
	NewARC(string, int) vela.Shared
	NewLFU(string, int) vela.Shared
}

var EmptySharedError = fmt.Errorf("share not found")

type EmptyShared struct {
	name string
	cap  int
}

func (esd *EmptyShared) Get(string) (interface{}, error)    { return nil, EmptySharedError }
func (esd *EmptyShared) Set(string, interface{}, int) error { return EmptySharedError }
func (esd *EmptyShared) Del(string)                         {}
func (esd *EmptyShared) Clear()                             {}

func (env *Environment) InitSharedEnv(v interface{}) {
	if env.shm != nil {
		env.Error("vela.shared hub already running")
		return
	}

	if shm, ok := v.(shared); ok {
		env.shm = shm
		return
	}
	env.Errorf("invalid shared object , got %p", v)
}

func (env *Environment) NewLRU(name string, cap int) vela.Shared {
	if env.shm == nil {
		return &EmptyShared{name, cap}
	}
	return env.shm.NewLRU(name, cap)
}

func (env *Environment) NewARC(name string, cap int) vela.Shared {
	if env.shm == nil {
		return &EmptyShared{name, cap}
	}
	return env.shm.NewARC(name, cap)
}

func (env *Environment) NewLFU(name string, cap int) vela.Shared {
	if env.shm == nil {
		return &EmptyShared{name, cap}
	}
	return env.shm.NewLFU(name, cap)
}
