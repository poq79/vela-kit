package variable

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"sync"
)

var xEnv vela.Environment

type Hub struct {
	mutex sync.RWMutex
	dict  map[string]variable
}

func (hub *Hub) String() string                         { return fmt.Sprintf("variable %p", hub) }
func (hub *Hub) Type() lua.LValueType                   { return lua.LTObject }
func (hub *Hub) AssertFloat64() (float64, bool)         { return 0, false }
func (hub *Hub) AssertString() (string, bool)           { return "", false }
func (hub *Hub) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (hub *Hub) Peek() lua.LValue                       { return hub }

func newVariableHub() *Hub {
	return &Hub{
		dict: make(map[string]variable, 16),
	}
}

func (hub *Hub) find(key string, entry *variable) bool {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	item, ok := hub.dict[key]
	if ok {
		*entry = item
		return ok
	}
	return false
}

func Constructor(env vela.Environment, callback func(*Hub) error) {
	xEnv = env
	hub := newVariableHub()
	hub.define(env)

	xEnv.Set("var", hub)
	xEnv.Set("readonly", lua.NewFunction(readOnlyItem))

	if e := callback(hub); e != nil {
		xEnv.Errorf("callback variable fail %v", e)
		return
	}
}
