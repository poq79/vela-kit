package audit

import (
	"errors"
	"fmt"
	auxlib "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"reflect"
)

func (ev *Event) String() string                         { return lua.B2S(ev.Byte()) }
func (ev *Event) Type() lua.LValueType                   { return lua.LTObject }
func (ev *Event) AssertFloat64() (float64, bool)         { return 0, false }
func (ev *Event) AssertString() (string, bool)           { return "", false }
func (ev *Event) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (ev *Event) Peek() lua.LValue                       { return ev }

func (ev *Event) ret(L *lua.LState) int {
	L.Push(ev)
	return 1
}

func (ev *Event) timeL(L *lua.LState) int {
	tm := L.CheckString(1)
	ev.time = auxlib.ToTime(tm)
	return ev.ret(L)
}

func (ev *Event) msgL(L *lua.LState) int {
	ev.msg = auxlib.Format(L, 0)
	return ev.ret(L)
}

func (ev *Event) subjectL(L *lua.LState) int {
	ev.subject = auxlib.Format(L, 0)
	return ev.ret(L)
}

func (ev *Event) fromL(L *lua.LState) int {
	lv := L.Get(1)
	switch lv.Type() {
	case lua.LTString:
		ev.from = lv.String()

	default:
		ev.from = L.CodeVM()
	}
	return ev.ret(L)
}

func (ev *Event) remoteL(L *lua.LState) int {
	ev.Remote(L.IsString(1))
	return ev.ret(L)
}

func (ev *Event) userL(L *lua.LState) int {
	ev.User(L.IsString(1))
	return ev.ret(L)
}

func (ev *Event) authL(L *lua.LState) int {
	ev.Auth(L.IsString(1))
	return ev.ret(L)
}

func (ev *Event) regionL(L *lua.LState) int {
	ev.region = L.IsString(1)
	return ev.ret(L)
}

func (ev *Event) portL(L *lua.LState) int {
	lv := L.Get(1)
	switch lv.Type() {
	case lua.LTNumber:
		ev.Port(int(lv.(lua.LNumber)))
	case lua.LTInt:
		ev.Port(int(lv.(lua.LInt)))
	default:
		ev.Port(auxlib.ToInt(lv.String()))
	}
	return ev.ret(L)
}

func (ev *Event) errL(L *lua.LState) int {
	val := L.Get(1)
	switch val.Type() {
	case lua.LTNil:
		//
	default:
		ev.E(errors.New(val.String()))
	}

	return ev.ret(L)
}

func (ev *Event) logL(L *lua.LState) int {
	if !L.IsFalse(1) {
		ev.Log()
	}

	return ev.ret(L)
}

func (ev *Event) putL(L *lua.LState) int {
	n := L.GetTop()

	if n >= 1 && L.IsTrue(1) {
		ev.Log()
	}

	if n >= 2 {
		ev.alert = L.IsTrue(2)
	}

	if n >= 3 {
		ev.Level(L.IsInt(3))
	}

	ev.Put()
	return ev.ret(L)
}

func (ev *Event) levelL(L *lua.LState) int {
	ev.Level(L.IsInt(1))
	return ev.ret(L)
}

func (ev *Event) alertL(L *lua.LState) int {
	ev.alert = L.IsFalse(1)
	return ev.ret(L)
}

func (ev *Event) Index(L *lua.LState, key string) lua.LValue {

	switch key {
	case "ID":
		return lua.S2L(ev.id)
	case "inet":
		return lua.S2L(ev.inet)
	case "time":
		v, e := ev.time.MarshalJSON()
		if e != nil {
			xEnv.Errorf("Event time to json error %v", e)
			return lua.LSNull
		}
		return lua.B2L(v)

	case "remote":
		return lua.S2L(ev.rAddr)

	case "port":
		return lua.LInt(ev.rPort)

	case "typeof":
		return lua.S2L(ev.typeof)

	case "auth":
		return lua.S2L(ev.auth)

	case "msg":
		return lua.S2L(ev.msg)

	case "subject":
		return lua.S2L(ev.subject)

	case "debug":
		return lua.B2L(ev.Byte())

	case "string":
		return lua.B2L(ev.Byte())

	case "from":
		return lua.S2L(ev.from)

	case "error":
		if ev.err == nil {
			return lua.LNil
		}

		return lua.S2L(ev.err.Error())

	case "region":
		return lua.S2L(ev.region)

	case "alert":
		return lua.LBool(ev.alert)

	case "Time":
		return L.NewFunction(ev.timeL)

	case "Msg":
		return L.NewFunction(ev.msgL)

	case "Subject":
		return L.NewFunction(ev.subjectL)

	case "From":

		return L.NewFunction(ev.fromL)

	case "Remote":
		return L.NewFunction(ev.remoteL)

	case "utime":
		return L.NewFunction(ev.userL)

	case "Auth":
		return L.NewFunction(ev.authL)

	case "Region":
		return L.NewFunction(ev.regionL)

	case "Port":
		return L.NewFunction(ev.portL)

	case "E":
		return L.NewFunction(ev.errL)

	case "Log":
		return L.NewFunction(ev.logL)

	case "Put":
		return L.NewFunction(ev.putL)

	case "Level":
		return L.NewFunction(ev.levelL)

	case "Alert":
		return L.NewFunction(ev.alertL)

	default:
		return lua.LNil

	}

}

func (ev *Event) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "err":
		ev.E(errors.New(val.String()))

	case "time":
		ev.Time(auxlib.ToTime(val.String()))

	case "remote":
		ev.Remote(val.String())

	case "port":
		switch val.Type() {
		case lua.LTNumber:
			ev.Port(int(val.(lua.LNumber)))
		case lua.LTInt:
			ev.Port(int(val.(lua.LInt)))
		default:
			ev.Port(auxlib.ToInt(val.String()))
		}

	case "region":
		ev.region = val.String()

	case "from":
		switch val.Type() {
		case lua.LTString:
			ev.from = val.String()
		default:
			ev.from = L.CodeVM()
		}

	case "msg":
		ev.msg = val.String()

	case "subject":
		ev.subject = val.String()
	case "user":
		ev.user = val.String()
	case "auth":
		ev.auth = val.String()
	case "level":
		ev.Level(lua.IsInt(val))
	case "typeof":
		ev.typeof = val.String()
	case "alert":
		ev.alert = lua.IsTrue(val)
	}

}

func CheckEvent(L *lua.LState, idx int) *Event {
	v := L.Get(idx)
	var ev *Event
	var ok bool

	switch v.Type() {
	case lua.LTObject:
		ev, ok = v.(*Event)

	case lua.LTAnyData:
		ev, ok = v.(*lua.AnyData).Data.(*Event)

	default:
		xEnv.Errorf("got lua %s , not lua vela-event", v.Type().String())
		return nil
	}

	if ok {
		return ev
	} else {
		xEnv.Errorf("not vela-event")
		return nil
	}
}

func brief(raw string) string {
	n := len(raw)
	if n < 100 {
		return raw
	} else {
		return raw[:100]
	}
}

func t2s(L *lua.LState) string {
	v := L.Get(1)
	var s string
	switch v.Type() {
	case lua.LTNil:
		s = fmt.Sprintf("lua:nil go:nil value:nil")
	case lua.LTInt:
		s = fmt.Sprintf("lua:int go:int value:%d", v.(lua.LInt))
	case lua.LTAnyData:
		obj := v.(*lua.AnyData).Data
		s = fmt.Sprintf("lua:any go: %s value:%v", reflect.TypeOf(obj).String(), obj)
	case lua.LTBool:
		s = fmt.Sprintf("lua:bool go:bool value:%v", v)
	case lua.LTFunction:
		s = fmt.Sprintf("lua:function go:func(*state) int")
	case lua.LTTable:
		s = fmt.Sprintf("lua:table")
	case lua.LTString:
		s = fmt.Sprintf("lua:string go:string value:%v", v)
	case lua.LTNumber:
		s = fmt.Sprintf("lua:number go:float64 value:%v", v)
	case lua.LTChannel:
		s = fmt.Sprintf("lua:channel go:channel value:%v", v)
	case lua.LTVelaData:
		obj := v.(*lua.VelaData).Data
		s = fmt.Sprintf("lua:lightUserData go:%s value:%v", reflect.TypeOf(obj).String(), obj)
	case lua.LTKv:
		s = fmt.Sprintf("lua:kv go:slice")
	case lua.LTThread:
		s = fmt.Sprintf("lua:thread go:*lua.LState")
	case lua.LTObject:
		s = fmt.Sprintf("lua:object go:%T value:%v", v, v)

	default:
		s = fmt.Sprintf("lua:object go:%T value:%v", v, v)
	}

	if L.Console != nil {
		L.Console.Println(s)
	}

	return brief(s)
}

func newLuaEvent(L *lua.LState) int {
	val := L.Get(1)
	ev := NewEvent("unknown").From(L.CodeVM())
	switch val.Type() {
	case lua.LTString:
		ev.typeof = val.String()

	case lua.LTTable:
		val.(*lua.LTable).Range(func(key string, val lua.LValue) {
			ev.NewIndex(L, key, val)
		})

	default:
		//todo
	}
	L.Push(ev)
	return 1
}

func newLuaDebug(L *lua.LState) int {
	line := L.Where(1)
	ev := NewEvent("logger").Subject("调试信息").From(L.CodeVM())
	ev.msg = "[" + line[:len(line)-1] + "] " + lua.Format(L, 0)

	ev.Log().Put()
	return 0
}

func newLuaError(L *lua.LState) int {
	line := L.Where(1)
	ev := NewEvent("logger").Subject("错误信息").From(L.CodeVM())
	ev.msg = line + lua.Format(L, 0)
	ev.High().Log().Put()
	return 0
}

func newLuaInfo(L *lua.LState) int {
	ev := NewEvent("logger").Subject("事件信息").From(L.CodeVM())
	ev.msg = auxlib.Format(L, 0)
	ev.Log().Put()
	return 0
}

func newLuaObjectType(L *lua.LState) int {
	ev := NewEvent("logger").Subject("类型分析").From(L.CodeVM())
	ev.msg = t2s(L)
	ev.Log().Put()
	return 0
}
