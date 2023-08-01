package shared

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"sync"
)

type shared struct {
	sync.RWMutex
	data map[string]*ShareBucket
}

func (s *shared) reset(name string) {
	if _, ok := s.data[name]; ok {
		//L.RaiseError("shared.%s exist" , key)
		xEnv.Errorf("%s shared clean", name)
	}
}

func (s *shared) NewLRU(name string, cap int) vela.Shared {
	bkt := New(cap).LRU().build()
	return newShareBucket(name, bkt)
}

func (s *shared) NewARC(name string, cap int) vela.Shared {
	bkt := New(cap).ARC().build()
	return newShareBucket(name, bkt)
}

func (s *shared) NewLFU(name string, cap int) vela.Shared {
	bkt := New(cap).LFU().build()
	return newShareBucket(name, bkt)
}
