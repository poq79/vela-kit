package audit

import (
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"reflect"
	"sync"
)

var (
	once sync.Once
	xEnv vela.Environment
)

var typeof = reflect.TypeOf((*Audit)(nil)).String()

type Audit struct {
	lua.SuperVelaData
	cfg *config
	fd  io.WriteCloser
}

func withConfig(cfg *config) *Audit {
	adt := &Audit{cfg: cfg}
	adt.V(lua.VTInit, typeof)
	return adt
}

func New() *Audit {
	return withConfig(velaMinConfig())
}

func (adt *Audit) openFile() {
	w := &lumberjack.Logger{
		Filename: adt.cfg.file,
		MaxSize:  1024 * 1024, //1G
		MaxAge:   180,
		Compress: false,
	}

	adt.fd = w
}

func (adt *Audit) output(ev *Event) {
	if adt.cfg.sdk != nil {
		adt.cfg.sdk.Write(ev.Byte())
	}

	if adt.fd != nil {
		adt.fd.Write(ev.Byte())
		adt.fd.Write([]byte("\n"))
	}
}

func (adt *Audit) pass(ev *Event) bool {
	n := len(adt.cfg.pass)
	if n == 0 {
		return false
	}

	for i := 0; i < n; i++ {
		if adt.cfg.pass[i](ev) {
			return true
		}
	}

	return false
}

func (adt *Audit) inhibit(ev *Event) {
	if !ev.alert {
		return
	}

	n := len(adt.cfg.rate)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		if adt.cfg.rate[i](adt.cfg.bkt, ev) {
			ev.alert = false
			return
		}
	}
}

func (adt *Audit) handle(ev *Event) {
	adt.output(ev)
	if adt.pass(ev) {
		xEnv.Debugf("by pass ev %s %s %s", ev.from, ev.typeof, ev.msg)
		return
	}

	//告警限速
	if ev.alert && !xEnv.IsDebug() {
		adt.inhibit(ev)
	}

	//流处理
	adt.cfg.pipe.Do(ev, adt.cfg.co, func(err error) {
		xEnv.Errorf("%v", err)
	})

	//是否上传
	if !ev.upload {
		return
	}

	err := xEnv.Push("/api/v1/broker/audit/event", json.RawMessage(ev.Byte()))
	if err != nil {
		//xEnv.Errorf("%s tnl send event fail %v", xEnv.TnlName(), err)
		return
	}
	//xEnv.Debugf("%s tnl send %v event succeed", xEnv.TnlName(), ev)
}
