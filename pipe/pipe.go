package pipe

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

type Fn func(...interface{}) error

type Chains struct {
	chain []Fn
	seek  int
	xEnv  vela.Environment
}

func (px *Chains) Merge(sub *Chains) {
	if sub.Len() == 0 {
		return
	}
	px.chain = append(px.chain, sub.chain...)
}

func (px *Chains) clone(co *lua.LState) *lua.LState {
	if px.xEnv == nil {
		px.xEnv = vela.GxEnv()
	}

	if co == nil {
		return px.xEnv.Coroutine()
	}

	return px.xEnv.Clone(co)
}
func (px *Chains) append(v Fn) {
	if v == nil {
		return
	}

	px.chain = append(px.chain, v)
}

func (px *Chains) coroutine() *lua.LState {
	if px.xEnv != nil {
		return px.xEnv.Coroutine()
	}
	return vela.GxEnv().Coroutine()
}

func (px *Chains) free(co *lua.LState) {
	if px.xEnv != nil {
		px.xEnv.Free(co)
		return
	}
	vela.GxEnv().Free(co)
}

func (px *Chains) invalid(format string, v ...interface{}) {
	if px.xEnv == nil {
		//vela.GxEnv().Errorf(format, v...)
		return
	}

	px.xEnv.Errorf(format, v...)
}

func New(opt ...func(*Chains)) (px *Chains) {
	px = &Chains{}

	for _, fn := range opt {
		fn(px)
	}

	return
}
