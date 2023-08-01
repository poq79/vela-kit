package shared

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

func (s *shared) String() string                         { return "vela.share.bucket" }
func (s *shared) Type() lua.LValueType                   { return lua.LTObject }
func (s *shared) AssertFloat64() (float64, bool)         { return 0, false }
func (s *shared) AssertString() (string, bool)           { return "", false }
func (s *shared) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (s *shared) Peek() lua.LValue                       { return s }

func (s *shared) NewIndex(L *lua.LState, key string, val lua.LValue) {

	if !L.CheckCodeVM("startup") {
		L.RaiseError("vela.shared.* not allow")
		return
	}

	if e := auxlib.Name(key); e != nil {
		L.RaiseError("%v", e)
		return
	}

	var bkt Cache
	var name string
	size := lua.CheckInt(L, val)

	//加锁
	s.Lock()
	defer s.Unlock()

	switch key[:3] {
	case "arc":
		name = key[4:]
		bkt = New(size).ARC().build()
	case "lru":
		name = key[4:]
		bkt = New(size).LRU().build()
	case "lfu":
		name = key[4:]
		bkt = New(size).LFU().build()
	default:
		name = key
		bkt = New(size).ARC().build()
	}

	s.reset(name)
	s.data[key] = newShareBucket(name, bkt)
}

func (s *shared) Index(L *lua.LState, key string) lua.LValue {
	s.RLock()
	defer s.RUnlock()

	if bkt, ok := s.data[key]; ok {
		return bkt
	}

	return lua.LNil
}
