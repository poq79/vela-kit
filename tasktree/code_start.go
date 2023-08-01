package tasktree

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"time"
)

type start struct {
	face lua.VelaEntry
	co   *lua.LState
	onE  func(error) //onError
	code string
}

func Start(co *lua.LState, face lua.VelaEntry) vela.Start {
	if co.CodeVM() == "" {
		co.RaiseError("not allow start vela in vm without code")
		return nil
	}
	return &start{face: face, co: co}

}

func (s *start) Err(fn func(error)) vela.Start {
	if fn != nil {
		s.onE = fn
	}
	return s
}

func (s *start) From(code string) vela.Start {
	s.code = code
	return s
}

func (s *start) run() error {
	switch s.face.State() {

	case lua.VTRun:
		obj, ok := s.face.(interface{ Reload() error })
		if ok {
			return obj.Reload()
		}

		if e := s.face.Close(); e != nil {
			return fmt.Errorf("%s close error %v", s.face.Name(), e)
		}

		return s.face.Start()
	default:
		return s.face.Start()
	}

}

func (s *start) Do() {
	if s.face == nil {
		return
	}

	if s.code != "" && s.co.CodeVM() != s.code {
		s.co.RaiseError("not allow with %s , must be %s", s.co.CodeVM(), s.code)
	}

	err := s.run()
	if err == nil {
		xEnv.Errorf("task.%s.%s start succeed", s.co.CodeVM(), s.face.Name())
		s.face.V(lua.VTRun, time.Now())
		return
	}

	if s.onE != nil {
		s.face.V(lua.VTErr, time.Now())
		s.onE(err)
	} else {
		s.face.V(lua.VTErr, time.Now())
		s.co.RaiseError("task.%s.%s start fail %v ", s.code, s.face.Name(), err)
	}

	if s.co != nil {
		s.co.RaiseError("%v", err)
	}
}
