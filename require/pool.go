package require

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"sync"
	"time"
)

type pool struct {
	mu   sync.RWMutex
	data []*cache
	stop chan struct{}
}

type Pool interface {
	RequireL(L *lua.LState, filename string) lua.LValue
	Require(filename string) lua.LValue
	Close() error
}

func newPool() *pool {
	return &pool{stop: make(chan struct{}, 1)}
}

func (p *pool) lock() {
	p.mu.Lock()
}

func (p *pool) unLock() {
	p.mu.Unlock()
}

func (p *pool) rLock() {
	p.mu.RLock()
}

func (p *pool) rUnlock() {
	p.mu.RUnlock()
}

func (p *pool) get(filename string) *cache {
	p.rLock()
	defer p.rUnlock()

	n := len(p.data)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		item := p.data[i]
		if item.name == filename {
			item.Hit()
			return item
		}
	}

	return nil
}

func (p *pool) append(item *cache) {
	p.lock()
	defer p.unLock()
	p.data = append(p.data, item)
}

func (p *pool) delete(idx int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	n := len(p.data)

	//item := p.data[idx]
	if idx == n-1 {
		p.data = p.data[:idx]
	} else {
		p.data = append(p.data[:idx], p.data[idx+1:]...)
	}

}

func (p *pool) flush() {
	n := len(p.data)
	if n == 0 {
		return
	}

	for i := 0; i < n; {
		item := p.data[i]
		mt := item.stat()

		if mt == item.mtime {
			i++
			continue
		}

		if mt == 0 {
			p.delete(i)
			item.gc()
			n -= 1 //遍历少一次
			xEnv.Errorf("%s 3rd delete succeed", item.file())

		} else {
			i++
			if err := item.load(); err != nil {
				xEnv.Errorf("%s 3rd update error %v", item.file(), err)
			} else {
				xEnv.Errorf("%s 3rd update succeed", item.file())

			}
		}
	}

}

func (p *pool) sync() {
	tk := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-tk.C:
			p.flush()
		case <-p.stop:
			return
		}
	}
}

func (p *pool) require(L *lua.LState, filename string) lua.LValue {
	var item *cache

	item = p.get(filename)
	if item != nil {
		return item
	}

	item = &cache{co: L, name: filename, status: OOP}

	err := item.load()
	if err != nil {
		xEnv.Errorf("compile 3rd fail %v", err)
		return lua.LNil
	}

	p.append(item)
	return item
}

func (p *pool) Close() error {
	p.stop <- struct{}{}
	return nil
}

func (p *pool) Name() string {
	return "require.pool"
}

func (p *pool) RequireL(L *lua.LState, filename string) lua.LValue {
	co := xEnv.Clone(L)

	return p.require(co, filename)
}

func (p *pool) Require(filename string) lua.LValue {
	L := xEnv.Coroutine()
	defer xEnv.Free(L)

	return p.require(L, filename)
}
