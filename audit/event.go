package audit

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"time"
)

type Event struct {
	time    time.Time //time
	id      string
	inet    string
	subject string //subject
	rAddr   string //remote addr
	rPort   int    //remote port
	from    string //Event from proc name
	typeof  string //type
	user    string //user
	auth    string //auth
	msg     string //info
	err     error  //error
	region  string
	alert   bool
	upload  bool
	level   string
}

func NewEvent(typeof string, opts ...func(*Event)) *Event {
	ev := &Event{
		id:     xEnv.ID(),
		inet:   xEnv.LocalAddr(),
		time:   time.Now(),
		level:  NOTICE,
		typeof: typeof,
	}

	n := len(opts)
	if n == 0 {
		return ev
	}

	for i := 0; i < n; i++ {
		opts[i](ev)
	}

	return ev
}

func Errorf(format string, v ...interface{}) *Event {
	return NewEvent("logger").Subject("发现错误").Msg(format, v...)
}

func Infof(format string, v ...interface{}) *Event {
	return NewEvent("logger").Subject("打印信息").Msg(format, v...)
}

func Debug(format string, v ...interface{}) *Event {
	return NewEvent("logger").Subject("调试信息").Msg(format, v...).From("vela-inline")
}

func Console(L *lua.LState, format string, v ...interface{}) *Event {
	line := L.Where(1)
	ev := NewEvent("logger").Subject("调试信息").From(L.CodeVM())
	ev.msg = "[" + line[:len(line)-1] + "] " + fmt.Sprintf(format, v...)
	return ev
}
