package lua

import (
	"context"
	"fmt"
	"os"
)

type LValueType int

const (
	LTNil LValueType = iota
	LTBool
	LTNumber
	LTInt
	LTUint
	LTInt64
	LTUint64
	LTString
	LTFunction
	LTUserData
	LTThread
	LTTable
	LTChannel
	LTVelaData
	LTSlice
	LTMap
	LTKv
	LTSkv
	LTAnyData
	LTObject
	LTGoFunction
	LTGoFuncErr
	LTGoFuncStr
	LTGoFuncInt
)

var lValueNames = [...]string{"nil", "boolean", "number", "int", "uint", "int64", "uint64", "string", "function", "userdata", "thread", "table", "channel", "veladata", "slice", "safe_map", "kv", "safe_kv", "AnyData", "object", "GoFunction", "GoFuncErr", "GoFuncStr", "GoFuncInt"}

func (vt LValueType) String() string {
	return lValueNames[int(vt)]
}

type LValue interface {
	String() string
	Type() LValueType
	// to reduce `runtime.assertI2T2` costs, this method should be used instead of the type assertion in heavy paths(typically inside the VM).
	AssertFloat64() (float64, bool)
	// to reduce `runtime.assertI2T2` costs, this method should be used instead of the type assertion in heavy paths(typically inside the VM).
	AssertString() (string, bool)
	// to reduce `runtime.assertI2T2` costs, this method should be used instead of the type assertion in heavy paths(typically inside the VM).
	AssertFunction() (*LFunction, bool)

	Peek() LValue
}

// LVIsFalse returns true if a given LValue is a nil or false otherwise false.
func LVIsFalse(v LValue) bool { return v == LNil || v == LFalse }

// LVIsFalse returns false if a given LValue is a nil or false otherwise true.
func LVAsBool(v LValue) bool { return v != LNil && v != LFalse }

// LVAsString returns string representation of a given LValue
// if the LValue is a string or number, otherwise an empty string.
func LVAsString(v LValue) string {
	switch sn := v.(type) {
	case LString, LNumber:
		return sn.String()
	default:
		return ""
	}
}

// LVCanConvToString returns true if a given LValue is a string or number
// otherwise false.
func LVCanConvToString(v LValue) bool {
	switch v.(type) {
	case LString, LNumber:
		return true
	default:
		return false
	}
}

// LVAsNumber tries to convert a given LValue to a number.
func LVAsNumber(v LValue) LNumber {
	switch lv := v.(type) {
	case LNumber:
		return lv
	case LString:
		if num, err := parseNumber(string(lv)); err == nil {
			return num
		}
	}
	return LNumber(0)
}

type LNilType struct{}

func (nl *LNilType) String() string                     { return "nil" }
func (nl *LNilType) Type() LValueType                   { return LTNil }
func (nl *LNilType) AssertFloat64() (float64, bool)     { return 0, false }
func (nl *LNilType) AssertString() (string, bool)       { return "", false }
func (nl *LNilType) AssertFunction() (*LFunction, bool) { return nil, false }
func (nl *LNilType) Peek() LValue                       { return nl }

var LNil = LValue(&LNilType{})

type LBool bool

func (bl LBool) String() string {
	if bool(bl) {
		return "true"
	}
	return "false"
}
func (bl LBool) Type() LValueType                   { return LTBool }
func (bl LBool) AssertFloat64() (float64, bool)     { return 0, false }
func (bl LBool) AssertString() (string, bool)       { return "", false }
func (bl LBool) AssertFunction() (*LFunction, bool) { return nil, false }
func (bl LBool) Peek() LValue                       { return bl }

var LTrue = LBool(true)
var LFalse = LBool(false)

type LString string

func (st LString) String() string                     { return string(st) }
func (st LString) Type() LValueType                   { return LTString }
func (st LString) AssertFloat64() (float64, bool)     { return 0, false }
func (st LString) AssertString() (string, bool)       { return string(st), true }
func (st LString) AssertFunction() (*LFunction, bool) { return nil, false }
func (st LString) Peek() LValue                       { return st }

// fmt.Formatter interface
func (st LString) Format(f fmt.State, c rune) {
	switch c {
	case 'd', 'i':
		if nm, err := parseNumber(string(st)); err != nil {
			defaultFormat(nm, f, 'd')
		} else {
			defaultFormat(string(st), f, 's')
		}
	default:
		defaultFormat(string(st), f, c)
	}
}

func (nm LNumber) String() string {
	if isInteger(nm) {
		return fmt.Sprint(int64(nm))
	}
	return fmt.Sprint(float64(nm))
}

func (nm LNumber) Type() LValueType                   { return LTNumber }
func (nm LNumber) AssertFloat64() (float64, bool)     { return float64(nm), true }
func (nm LNumber) AssertString() (string, bool)       { return "", false }
func (nm LNumber) AssertFunction() (*LFunction, bool) { return nil, false }
func (nm LNumber) Peek() LValue                       { return nm }

// fmt.Formatter interface
func (nm LNumber) Format(f fmt.State, c rune) {
	switch c {
	case 'q', 's':
		defaultFormat(nm.String(), f, c)
	case 'b', 'c', 'd', 'o', 'x', 'X', 'U':
		defaultFormat(int64(nm), f, c)
	case 'e', 'E', 'f', 'F', 'g', 'G':
		defaultFormat(float64(nm), f, c)
	case 'i':
		defaultFormat(int64(nm), f, 'd')
	default:
		if isInteger(nm) {
			defaultFormat(int64(nm), f, c)
		} else {
			defaultFormat(float64(nm), f, c)
		}
	}
}

type LTable struct {
	Metatable LValue

	array   []LValue
	dict    map[LValue]LValue
	strdict map[string]LValue
	keys    []LValue
	k2i     map[LValue]int
}

func (tb *LTable) String() string                     { return fmt.Sprintf("table: %p", tb) }
func (tb *LTable) Type() LValueType                   { return LTTable }
func (tb *LTable) AssertFloat64() (float64, bool)     { return 0, false }
func (tb *LTable) AssertString() (string, bool)       { return "", false }
func (tb *LTable) AssertFunction() (*LFunction, bool) { return nil, false }
func (tb *LTable) Peek() LValue                       { return tb }

type LFunction struct {
	IsG       bool
	Env       *LTable
	Proto     *FunctionProto
	GFunction LGFunction
	Upvalues  []*Upvalue
}
type LGFunction func(*LState) int

func (fn *LFunction) String() string                     { return fmt.Sprintf("function: %p", fn) }
func (fn *LFunction) Type() LValueType                   { return LTFunction }
func (fn *LFunction) AssertFloat64() (float64, bool)     { return 0, false }
func (fn *LFunction) AssertString() (string, bool)       { return "", false }
func (fn *LFunction) AssertFunction() (*LFunction, bool) { return fn, true }
func (fn *LFunction) Peek() LValue                       { return fn }

type Global struct {
	MainThread    *LState
	CurrentThread *LState
	Registry      *LTable
	Global        *LTable

	builtinMts map[int]LValue
	tempFiles  []*os.File
	gccount    int32
}

type LState struct {
	G       *Global
	Parent  *LState
	Env     *LTable
	Panic   func(*LState)
	Dead    bool
	Options Options

	stop         int32
	reg          *registry
	stack        callFrameStack
	alloc        *allocator
	currentFrame *callFrame
	wrapped      bool
	uvcache      *Upvalue
	hasErrorFunc bool
	mainLoop     func(*LState, *callFrame)
	ctx          context.Context
	ctxCancelFn  context.CancelFunc
	Exdata       interface{}
	Console      Console
	metadata     [7]interface{}
}

func (ls *LState) String() string                     { return fmt.Sprintf("thread: %p", ls) }
func (ls *LState) Type() LValueType                   { return LTThread }
func (ls *LState) AssertFloat64() (float64, bool)     { return 0, false }
func (ls *LState) AssertString() (string, bool)       { return "", false }
func (ls *LState) AssertFunction() (*LFunction, bool) { return nil, false }
func (ls *LState) Peek() LValue                       { return ls }

type LUserData struct {
	Value     interface{}
	Env       *LTable
	Metatable LValue
}

func (ud *LUserData) String() string                     { return fmt.Sprintf("userdata: %p", ud) }
func (ud *LUserData) Type() LValueType                   { return LTUserData }
func (ud *LUserData) AssertFloat64() (float64, bool)     { return 0, false }
func (ud *LUserData) AssertString() (string, bool)       { return "", false }
func (ud *LUserData) AssertFunction() (*LFunction, bool) { return nil, false }
func (ud *LUserData) Peek() LValue                       { return ud }

type LChannel chan LValue

func (ch LChannel) String() string                     { return fmt.Sprintf("channel: %p", ch) }
func (ch LChannel) Type() LValueType                   { return LTChannel }
func (ch LChannel) AssertFloat64() (float64, bool)     { return 0, false }
func (ch LChannel) AssertString() (string, bool)       { return "", false }
func (ch LChannel) AssertFunction() (*LFunction, bool) { return nil, false }
func (ch LChannel) Peek() LValue                       { return ch }
