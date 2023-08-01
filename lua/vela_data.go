package lua

import (
	"fmt"
)

const (
	VTInit VelaState = iota
	VTRun
	VTErr
	VTClose
	VTPanic
	VTPrivate
	VTMode
)

type VelaState uint32

var VelaStateValue = [...]string{"init", "run", "error", "close", "panic", "private", "mode"}

func (v VelaState) String() string {
	return VelaStateValue[int(v)]
}

type VelaEntry interface {
	Name() string     //获取当前对象名称
	Type() string     //获取对象类型
	State() VelaState //获取状态
	Start() error
	Close() error
	NewMeta(*LState, LValue, LValue)  //设置字段
	Meta(*LState, LValue) LValue      //获取字段
	Index(*LState, string) LValue     //获取字符串字段 __index function
	NewIndex(*LState, string, LValue) //设置字段 __newindex function

	Show(Console) //控制台打印
	Help(Console) //控制台 辅助信息
	//V  设置信息
	V(...interface{})
}

type VelaData struct {
	private bool
	code    string
	Data    VelaEntry
}

func IsChar(ch byte) bool {
	if ch >= 'a' && ch <= 'z' {
		return true
	}

	if ch >= 'A' && ch <= 'Z' {
		return true
	}

	return false
}

func IsIntChar(ch byte) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

func VelaNameE(v string) error {
	if len(v) < 2 {
		return InvalidVelaName
	}

	if !IsChar(v[0]) {
		return InvalidVelaName
	}

	n := len(v)
	for i := 1; i < n; i++ {
		ch := v[i]
		switch {
		case IsChar(ch), IsIntChar(ch):
			continue

		case ch == '_':
			continue

		default:
			return InvalidVelaName
		}
	}

	return nil
}

func NewVelaData(v VelaEntry) *VelaData {
	return &VelaData{Data: v, private: false}
}

func (vd *VelaData) String() string                     { return fmt.Sprintf("veladata: %p", vd) }
func (vd *VelaData) Type() LValueType                   { return LTVelaData }
func (vd *VelaData) AssertFloat64() (float64, bool)     { return 0, false }
func (vd *VelaData) AssertString() (string, bool)       { return "", false }
func (vd *VelaData) AssertFunction() (*LFunction, bool) { return nil, false }
func (vd *VelaData) Peek() LValue                       { return vd }

func (vd *VelaData) Close() error {
	if vd.Data == nil {
		return fmt.Errorf("%s not found entry", vd.code)
	}

	return vd.Data.Close()
}

func (vd *VelaData) Private(L *LState) {
	if !L.CheckCodeVM(vd.CodeVM()) {
		L.RaiseError("proc private with %s not allow, must be %s", L.CodeVM(), vd.CodeVM())
		return
	}
	vd.private = true
}

func (vd *VelaData) IsPrivate() bool {
	return vd.private
}

func (vd *VelaData) CodeVM() string {
	return vd.code
}

func (vd *VelaData) IsNil() bool {
	return vd.Data == nil
}

func (vd *VelaData) Set(v VelaEntry) {
	vd.Data = v
}
