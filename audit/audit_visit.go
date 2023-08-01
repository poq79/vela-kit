package audit

import (
	"errors"
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"time"
)

func (adt *Audit) Name() string {
	return adt.cfg.name
}

func (adt *Audit) E(err error) {
	if err == nil {
		return
	}

	xEnv.Errorf("vela audit handle fail , error: %v", err)
}

func (adt *Audit) Close() error {
	if !adt.IsRun() {
		return errors.New(adt.Name() + "can't close , err: is close")
	}

	adt.V(lua.VTClose)
	if adt.fd != nil {
		adt.fd.Close()
	}

	adt.cfg = velaMinConfig()
	return nil
}

func (adt *Audit) Start() error {
	if adt.IsRun() {
		return fmt.Errorf("%s is running", adt.Name())
	}
	adt.openFile()
	adt.V(lua.VTRun, time.Now())
	return nil
}

func (adt *Audit) Errorf(format string, v ...interface{}) *Event {
	return NewEvent("logger").Subject("发现错误").Msg(format, v...)
}
