package tasktree

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

func newCodeEv(c *Code, sub string) *audit.Event {
	return audit.NewEvent("task").Subject(sub).From(c.Key()).Alert()
}

func in(k string, s []string) bool {
	if len(s) == 0 {
		return false
	}

	for _, item := range s {
		if k == item {
			return true
		}
	}
	return false
}

func CheckCodeVM(L *lua.LState) (*Code, bool) {
	cname := L.CodeVM()
	if cname == "" {
		return nil, false
	}

	cd, ok := L.Exdata.(*Code)
	return cd, ok
}

func wakeup(co *lua.LState, way vela.Way) error {
	code, ok := CheckCodeVM(co)
	if !ok {
		return fmt.Errorf("invalid code object")
	}

	code.wakeup(co, way)

	return code.Wrap()
}

func wrap(co *lua.LState) error {
	code, ok := CheckCodeVM(co)
	if ok {
		return code.Wrap()
	}
	return nil
}
