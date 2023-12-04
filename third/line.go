package third

import (
	"bufio"
	"context"
	"fmt"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"os"
)

type Scanner struct {
	file string
	s    *bufio.Scanner
}

func (s *Scanner) String() string                         { return fmt.Sprintf("%p", s) }
func (s *Scanner) Type() lua.LValueType                   { return lua.LTObject }
func (s *Scanner) AssertFloat64() (float64, bool)         { return 0, false }
func (s *Scanner) AssertString() (string, bool)           { return "", false }
func (s *Scanner) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (s *Scanner) Peek() lua.LValue                       { return s }

func (s *Scanner) textL(L *lua.LState) int {
	chains := pipe.NewByLua(L, pipe.Env(xEnv))
	if chains.Len() <= 0 {
		return 0
	}

	ctx, cancel := context.WithCancel(L.Context())
	defer cancel()

	f, err := os.Open(s.file)
	if err != nil {
		L.Pushf("%v", err)
		return 1
	}
	defer f.Close()

	co := xEnv.Clone(L)
	defer xEnv.Free(co)

	sc := bufio.NewScanner(f)
	for sc.Scan() {

		select {
		case <-ctx.Done():
			return 0
		default:
			chains.Do(lua.S2L(sc.Text()), co, func(err error) {
				//todo
			})
		}
	}
	return 0
}

func (s *Scanner) jsonL(L *lua.LState) int {
	chains := pipe.NewByLua(L, pipe.Env(xEnv))
	if chains.Len() <= 0 {
		return 0
	}

	ctx, cancel := context.WithCancel(L.Context())
	defer cancel()

	f, err := os.Open(s.file)
	if err != nil {
		L.Pushf("%v", err)
		return 1
	}

	defer f.Close()
	co := xEnv.Clone(L)
	defer xEnv.Free(co)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		select {
		case <-ctx.Done():
			return 0
		default:
			val := sc.Bytes()
			tab, e := kind.Decode(L, val)
			if e != nil {
				//todo
			}
			chains.Do(tab, co, func(err error) {
				//todo
			})
		}
	}
	return 0
}

func (s *Scanner) tableL(L *lua.LState) int {
	ctx, cancel := context.WithCancel(L.Context())
	defer cancel()

	f, err := os.Open(s.file)
	if err != nil {
		L.Pushf("%v", err)
		return 1
	}
	defer f.Close()

	tab := L.NewTable()
	idx := 0

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		select {
		case <-ctx.Done():
			return 0
		default:
			idx++
			val := sc.Bytes()
			tab.RawSetInt(idx, lua.B2L(val))
		}
	}
	L.Push(tab)
	return 1
}

func (s *Scanner) Index(L *lua.LState, key string) lua.LValue {
	switch key {

	case "text":
		return lua.NewFunction(s.textL)
	case "json":
		return lua.NewFunction(s.jsonL)
	case "table":
		return lua.NewFunction(s.tableL)

	default:
		return lua.LNil
	}
}
