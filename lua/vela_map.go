package lua

import (
	"github.com/vela-ssoc/vela-kit/grep"
	"sync"
)

type Map struct {
	safe  bool
	mutex sync.RWMutex
	entry map[string]LValue
}

func NewMap(cap int, safe bool) *Map {
	return &Map{
		safe:  safe,
		entry: make(map[string]LValue, cap),
	}
}

func (m *Map) String() string                     { return B2S(m.Byte()) }
func (m *Map) Type() LValueType                   { return LTMap }
func (m *Map) AssertFloat64() (float64, bool)     { return 0, false }
func (m *Map) AssertString() (string, bool)       { return "", false }
func (m *Map) AssertFunction() (*LFunction, bool) { return nil, false }
func (m *Map) Peek() LValue                       { return m }

func (m *Map) Byte() []byte {
	if m.safe {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
	}

	enc := Json(m.Len())
	enc.Tab("")
	for k, v := range m.entry {
		if v == nil {
			enc.KV(k, nil)
			continue
		}

		switch v.Type() {
		case LTInt:
			enc.KV(k, int(v.(LInt)))
		case LTNumber:
			enc.KV(k, float64(v.(LNumber)))
		case LTBool:
			enc.KV(k, bool(v.(LBool)) == true)
		default:
			enc.KV(k, v.String())
		}
	}
	enc.End("}")
	return enc.Bytes()
}

func (m *Map) Keys() []string {
	if m.safe {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
	}

	keys := make([]string, len(m.entry))
	i := 0
	for key := range m.entry {
		keys[i] = key
		i++
	}

	return keys
}

func (m *Map) Len() int {
	if m.safe {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
	}
	return len(m.entry)
}

func (m *Map) Get(key string) (LValue, bool) {
	if m.safe {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
	}

	v, ok := m.entry[key]
	return v, ok
}

func (m *Map) Set(key string, val LValue) {
	if m.safe {
		m.mutex.Lock()
		defer m.mutex.Unlock()
	}

	m.entry[key] = val
}

func (m *Map) Index(L *LState, key string) LValue {
	v, ok := m.Get(key)
	if ok {
		return v
	}

	return LNil
}

func (m *Map) NewIndex(L *LState, key string, val LValue) {
	m.Set(key, val)
}

func (m *Map) Meta(L *LState, key LValue) LValue {
	k, ok := key.AssertString()
	if ok {
		return m.Index(L, k)
	}

	return LNil
}

func (m *Map) NewMeta(L *LState, key LValue, val LValue) {
	k, ok := key.AssertString()
	if ok {
		m.NewIndex(L, k, val)
		return
	}
}

func (m *Map) countL(L *LState) int {
	L.Push(LInt(len(m.entry)))
	return 1
}

func (m *Map) keyL(L *LState) int {
	filter := grep.New(L.IsString(1))

	if m.safe {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
	}

	n := len(m.entry)
	keys := make([]LValue, n)
	i := 0

	for key := range m.entry {
		if filter(key) {
			keys[i] = S2L(key)
			i++
		}
	}
	L.PushAny(keys[:i])
	return 1
}

func (m *Map) MetaTable(L *LState, key string) LValue {
	switch key {
	case "count":
		return NewFunction(m.countL)

	case "key":
		return NewFunction(m.keyL)

	}

	return LNil
}
