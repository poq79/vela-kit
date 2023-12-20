package lua

import (
	"fmt"
	"os"
	"strings"
)

const (
	M0 int = iota
	M1
	M2
	M3
	M4
	M5
	M6
)

func (ls *LState) CheckObject(n int) LValue {
	lv := ls.Get(n)

	if lv.Type() != LTObject {
		ls.TypeError(n, LTObject)
		return nil
	}
	return lv
}

func (ls *LState) PushAny(v interface{}) {
	ls.Push(ToLValue(v))
}

func (ls *LState) Pushf(format string, v ...interface{}) {
	ls.Push(LString(fmt.Sprintf(format, v...)))
}

func (ls *LState) CheckSocket(n int) string {
	v := ls.CheckString(n)
	if e := CheckSocket(v); e != nil {
		ls.RaiseError("must be socket , got fail , error:%v", e)
		return ""
	}
	return v

}

func (ls *LState) CheckSockets(n int) string {
	v := ls.CheckString(n)
	arr := strings.Split(v, ",")

	var err error
	for _, item := range arr {
		err = CheckSocket(item)
		if err != nil {
			ls.RaiseError("%s error: %v", err)
			return ""
		}
	}

	return v
}

func (ls *LState) CheckFile(n int) string {
	v := ls.CheckString(n)

	_, err := os.Stat(v)
	if os.IsNotExist(err) {
		ls.RaiseError("not found %s file", v)
		return ""
	}

	return v
}

func (ls *LState) IsTrue(n int) bool {
	return IsTrue(ls.Get(n))
}

func (ls *LState) IsFalse(n int) bool {
	return IsFalse(ls.Get(n))
}

func (ls *LState) IsNumber(n int) LNumber {
	return IsNumber(ls.Get(n))
}

func (ls *LState) IsInt(n int) int {
	return IsInt(ls.Get(n))
}

func (ls *LState) IsFunc(n int) *LFunction {
	return IsFunc(ls.Get(n))
}

func (ls *LState) IsString(n int) string {
	return IsString(ls.Get(n))
}

type CallBackFunction func(LValue) (stop bool)

func (ls *LState) Callback(fn CallBackFunction) {
	n := ls.GetTop()
	if n == 0 {
		return
	}

	for i := 1; i <= n; i++ {
		if fn(ls.Get(i)) {
			return
		}
	}
}

func (ls *LState) SetMetadata(id int, v interface{}) {
	switch id {
	case M0, M1, M2, M3, M4, M5:
		ls.metadata[id] = v
	case M6:
		if ls.metadata[id] != nil {
			ls.RaiseError("The storage with metadata of 7 is not empty")
			return
		}
		ls.metadata[id] = v
	default:
		ls.RaiseError("metadata over flow size:%d  id:%d", 7, id)
	}
}

func (ls *LState) Metadata(id int) interface{} {
	switch id {
	case M0, M1, M2, M3, M4, M5, M6:
		return ls.metadata[id]
	default:
		ls.RaiseError("metadata over flow size:%d  id:%d", 7, id)
		return nil
	}
}

func (ls *LState) A() interface{} {
	return ls.Metadata(0)
}

func (ls *LState) B() interface{} {
	return ls.Metadata(1)
}

func (ls *LState) C() interface{} {
	return ls.Metadata(2)
}

func (ls *LState) D() interface{} {
	return ls.Metadata(3)
}

func (ls *LState) E() interface{} {
	return ls.Metadata(4)
}

func (ls *LState) SetA(v interface{}) {
	ls.SetMetadata(0, v)
}

func (ls *LState) SetB(v interface{}) {
	ls.SetMetadata(1, v)
}
func (ls *LState) SetC(v interface{}) {
	ls.SetMetadata(2, v)
}
func (ls *LState) SetD(v interface{}) {
	ls.SetMetadata(3, v)
}
func (ls *LState) SetE(v interface{}) {
	ls.SetMetadata(4, v)
}

func (ls *LState) Copy(L *LState) {
	ls.Exdata = L.Exdata
	ls.ctx = L.ctx
	ls.Console = L.Console
	ls.metadata = L.metadata //数组深拷贝
}

func (ls *LState) Keepalive() {
	ls.Console = nil
	ls.Exdata = nil
	ls.SetContext(nil)
	ls.metadata = [7]interface{}{}
	ls.SetTop(0)
}

func (ls *LState) StackTrace(level int) string {
	return ls.stackTrace(level)
}

func (ls *LState) Use(value LValue) {
	if ls.Console != nil {
		ls.Console.Println(value.String())
	}

	ls.Push(value)
}
