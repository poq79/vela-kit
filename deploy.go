package vkit

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/agent"
	"github.com/vela-ssoc/vela-kit/env"
	"github.com/vela-ssoc/vela-kit/vela"
	"os"
)

type Deploy struct {
	name string
	all  bool
	use  func(env vela.Environment)
	doc  string
}

type EngineFunc func(*Deploy)

func All() EngineFunc {
	return func(e *Deploy) {
		e.all = true
	}
}

func Doc(doc string) EngineFunc {
	return func(dly *Deploy) {
		dly.doc = doc
	}
}

func Use(fn func(vela.Environment)) EngineFunc {
	return func(e *Deploy) {
		e.use = fn
	}
}

func New(name string, options ...EngineFunc) *Deploy {
	e := &Deploy{name: name}
	for _, fn := range options {
		fn(e)
	}
	return e
}

func (dly *Deploy) Agent() {
	if os.Args[1] == "version" {
		fmt.Println(dly.doc)
		return
	}
	agent.By(dly.name, dly.define())
}

func (dly *Deploy) Debug(hide Hide) {
	xEnv := env.Create("debug", dly.name, hide.Protect)
	dly.define()(xEnv)

	xEnv.Error("ssc sensor debug start")
	xEnv.Spawn(0, func() {
		xEnv.Dev(hide.Lan[0], hide.Vip[0], hide.Edition, hide.Hostname)
	})

	xEnv.Error("ssc sensor debug succeed")
	xEnv.Notify()
	xEnv.Error("ssc sensor exit succeed")
}
