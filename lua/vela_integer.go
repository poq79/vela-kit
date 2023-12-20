package lua

import (
	"github.com/spf13/cast"
	"strconv"
)

type LInt int
type LUint uint
type LInt64 int64
type LUint64 uint64

func (i LInt) Type() LValueType                   { return LTInt }
func (i LInt) String() string                     { return strconv.Itoa(int(i)) }
func (i LInt) AssertFloat64() (float64, bool)     { return float64(i), true }
func (i LInt) AssertString() (string, bool)       { return "", false }
func (i LInt) AssertFunction() (*LFunction, bool) { return nil, false }
func (i LInt) Peek() LValue                       { return i }

func (ui LUint) Type() LValueType                   { return LTUint }
func (ui LUint) String() string                     { return cast.ToString(ui) }
func (ui LUint) AssertFloat64() (float64, bool)     { return cast.ToFloat64(uint(ui)), true }
func (ui LUint) AssertString() (string, bool)       { return "", false }
func (ui LUint) AssertFunction() (*LFunction, bool) { return nil, false }
func (ui LUint) Peek() LValue                       { return ui }

func (i LInt64) Type() LValueType                   { return LTInt64 }
func (i LInt64) String() string                     { return strconv.Itoa(int(i)) }
func (i LInt64) AssertFloat64() (float64, bool)     { return float64(i), true }
func (i LInt64) AssertString() (string, bool)       { return "", false }
func (i LInt64) AssertFunction() (*LFunction, bool) { return nil, false }
func (i LInt64) Peek() LValue                       { return i }

func (ui LUint64) Type() LValueType                   { return LTUint64 }
func (ui LUint64) String() string                     { return strconv.FormatUint(uint64(ui), 10) }
func (ui LUint64) AssertFloat64() (float64, bool)     { return float64(ui), true }
func (ui LUint64) AssertString() (string, bool)       { return "", false }
func (ui LUint64) AssertFunction() (*LFunction, bool) { return nil, false }
func (ui LUint64) Peek() LValue                       { return ui }
