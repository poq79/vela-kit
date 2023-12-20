package lua

import "fmt"

type GoFunction func() error

func (fn GoFunction) String() string                     { return fmt.Sprintf("Go function error: %p", fn) }
func (fn GoFunction) Type() LValueType                   { return LTGoFunction }
func (fn GoFunction) AssertFloat64() (float64, bool)     { return 0, false }
func (fn GoFunction) AssertString() (string, bool)       { return "", false }
func (fn GoFunction) AssertFunction() (*LFunction, bool) { return nil, false }
func (fn GoFunction) Peek() LValue                       { return fn }

type GoFuncErr func(...interface{}) error

func (ge GoFuncErr) String() string                     { return fmt.Sprintf("Go function error: %p", ge) }
func (ge GoFuncErr) Type() LValueType                   { return LTGoFuncErr }
func (ge GoFuncErr) AssertFloat64() (float64, bool)     { return 0, false }
func (ge GoFuncErr) AssertString() (string, bool)       { return "", false }
func (ge GoFuncErr) AssertFunction() (*LFunction, bool) { return nil, false }
func (ge GoFuncErr) Peek() LValue                       { return ge }

type GoFuncStr func(...interface{}) string

func (gs GoFuncStr) String() string                     { return fmt.Sprintf("Go fucntion string: %p", gs) }
func (gs GoFuncStr) Type() LValueType                   { return LTGoFuncStr }
func (gs GoFuncStr) AssertFloat64() (float64, bool)     { return 0, false }
func (gs GoFuncStr) AssertString() (string, bool)       { return "", false }
func (gs GoFuncStr) AssertFunction() (*LFunction, bool) { return nil, false }
func (gs GoFuncStr) Peek() LValue                       { return gs }

type GoFuncInt func(...interface{}) int

func (gi GoFuncInt) String() string                     { return fmt.Sprintf("Go function int : %p", gi) }
func (gi GoFuncInt) Type() LValueType                   { return LTGoFuncInt }
func (gi GoFuncInt) AssertFloat64() (float64, bool)     { return 0, false }
func (gi GoFuncInt) AssertString() (string, bool)       { return "", false }
func (gi GoFuncInt) AssertFunction() (*LFunction, bool) { return nil, false }
func (gi GoFuncInt) Peek() LValue                       { return gi }
