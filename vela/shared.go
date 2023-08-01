package vela

type Shared interface {
	Get(string) (interface{}, error)
	Set(string, interface{}, int) error
	Del(string)
	Clear()
}

type SharedEnv interface {
	InitSharedEnv(interface{})
	NewLRU(string, int) Shared
	NewARC(string, int) Shared
	NewLFU(string, int) Shared
}
