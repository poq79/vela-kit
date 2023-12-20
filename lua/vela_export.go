package lua

type Export struct {
	name     string
	fn       *LFunction
	index    func(*LState, string) LValue
	newIndex func(*LState, string, LValue)
	tuple    UserKV
}

func (e Export) String() string                 { return e.name }
func (e Export) Type() LValueType               { return LTObject }
func (e Export) AssertFloat64() (float64, bool) { return 0, false }
func (e Export) AssertString() (string, bool)   { return "", false }
func (e Export) Peek() LValue                   { return e }

func (e Export) AssertFunction() (*LFunction, bool) {
	if e.fn == nil {
		return nil, false
	}
	return e.fn, true
}

func (e Export) NewIndex(L *LState, key string, val LValue) {
	if e.newIndex == nil {
		return
	}
	e.newIndex(L, key, val)
}

func (e Export) Index(L *LState, key string) LValue {
	var lv LValue

	if e.tuple == nil {
		goto index
	}

	lv = e.tuple.Get(key)
	if lv.Type() != LTNil {
		return lv
	}

index:
	if e.index == nil {
		return LNil
	}

	return e.index(L, key)
}

func NewExport(name string, opt ...func(*Export)) Export {
	e := Export{name: name}

	for _, fn := range opt {
		fn(&e)
	}
	return e
}

func WithFunc(gn LGFunction) func(e *Export) {
	return func(e *Export) {
		e.fn = NewFunction(gn)
	}
}

func WithIndex(ex func(*LState, string) LValue) func(e *Export) {
	return func(e *Export) {
		e.index = ex
	}
}

func WithNewIndex(nex func(*LState, string, LValue)) func(e *Export) {
	return func(e *Export) {
		e.newIndex = nex
	}
}

func WithTable(kv UserKV) func(e *Export) {
	return func(e *Export) {
		e.tuple = kv
	}
}
