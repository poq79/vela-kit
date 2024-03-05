package audit

import (
	"errors"
	"fmt"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"go.uber.org/zap/zapcore"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	DISASTER string = "紧急"
	HIGH     string = "重要"
	MIDDLE   string = "次要"
	NOTICE   string = "普通"
)

func (ev *Event) doRegion() {
	ip, err := xEnv.Region(ev.rAddr)
	if err != nil {
		xEnv.Debugf("vela-event region %s error %v", lua.B2L(ev.Byte()), err)
		return
	}
	ev.region = lua.B2S(ip.Byte())
}

func (ev *Event) Byte() []byte {
	if ev == nil {
		return []byte{}
	}

	buf := kind.NewJsonEncoder()
	buf.Tab("")
	buf.KV("time", ev.time)
	buf.KV("minion_id", ev.id)
	buf.KV("inet", ev.inet)
	buf.KV("subject", ev.subject)
	buf.KV("remote_addr", ev.rAddr)
	buf.KI("remote_port", ev.rPort)
	buf.KV("region", ev.region)
	buf.KV("from_code", ev.from)
	buf.KV("typeof", ev.typeof)
	buf.KV("user", ev.user)
	buf.KV("auth", ev.auth)
	buf.KV("msg", ev.msg)
	buf.KV("error", ev.err)
	buf.KV("alert", ev.alert)
	buf.KV("level", ev.level)
	buf.End("}")
	return buf.Bytes()
}

func (ev *Event) toLine() string {
	return fmt.Sprintf("[%s] [%s] [%s] %s  %s  %s  %d  %s  %s %s  %s  %s  %s  %v",
		ev.level, ev.time, ev.id, ev.inet, ev.subject, ev.rAddr, ev.rPort, ev.from, ev.typeof,
		ev.user, ev.auth, ev.msg, ev.region, ev.err)
}

func (ev *Event) json(L *lua.LState) int {
	L.Push(lua.B2L(ev.Byte()))
	return 1
}

func (ev *Event) line(L *lua.LState) int {
	L.Push(lua.S2L(ev.toLine()))
	return 1
}

func (ev *Event) ToValue(L *lua.LState) lua.LValue {
	return L.NewAnyData(ev)
}

func (ev *Event) Subject(format string, args ...interface{}) *Event {
	if len(args) == 0 {
		ev.subject = format
	} else {
		ev.subject = fmt.Sprintf(format, args)
	}

	return ev
}

func (ev *Event) Time(v time.Time) *Event {
	ev.time = v
	return ev
}

func (ev *Event) Remote(v interface{}) *Event {
	switch addr := v.(type) {
	case string:
		ip := addr
		port := -1

		s := strings.Split(addr, ":")
		if len(s) == 2 {
			ip = s[0]
			if p, err := strconv.Atoi(s[1]); err == nil {
				port = p
			}
		}

		if net.ParseIP(ip) != nil {
			ev.rAddr = ip
			ev.doRegion()
		}

		if port < 1 || port > 65535 {
			return ev
		}
		ev.rPort = port

	case net.IPNet:
		ev.rAddr = addr.IP.String()
		ev.doRegion()
	case net.IPAddr:
		ev.rAddr = addr.IP.String()
		ev.doRegion()

	case net.Conn:
		switch nt := addr.RemoteAddr().(type) {
		case *net.UDPAddr:
			ev.rAddr = nt.IP.String()
			ev.rPort = nt.Port
			ev.doRegion()

		case *net.TCPAddr:
			ev.rAddr = nt.IP.String()
			ev.rPort = nt.Port
			ev.doRegion()
		default:
			return ev
		}

	case net.Addr:
		switch nt := addr.(type) {
		case *net.UDPAddr:
			ev.rAddr = nt.IP.String()
			ev.rPort = nt.Port
			ev.doRegion()

		case *net.TCPAddr:
			ev.rAddr = nt.IP.String()
			ev.rPort = nt.Port
			ev.doRegion()
		default:
			return ev
		}
	}

	return ev
}

func (ev *Event) Port(v int) *Event {
	ev.rPort = v
	return ev
}

func (ev *Event) Task() string {
	return ev.from
}

func (ev *Event) Typeof() string {
	return ev.typeof
}

func (ev *Event) From(v string) *Event {
	ev.from = v
	return ev
}

func (ev *Event) Msg(format string, args ...interface{}) *Event {
	ev.msg = fmt.Sprintf(format, args...)
	return ev
}

func (ev *Event) User(v string) *Event {
	ev.user = v
	return ev
}

func (ev *Event) Auth(v string) *Event {
	ev.auth = v
	return ev
}

func (ev *Event) E(e error) *Event {
	ev.err = e
	return ev
}

func (ev *Event) Log() *Event {

	if ev.err == nil {
		xEnv.Debugf("[%s] [%s] %s %s %s %s %s %s %d %s",
			ev.level, ev.subject, ev.from, ev.typeof,
			ev.user, ev.auth, ev.msg, ev.rAddr, ev.rPort, ev.region)
		//xEnv.Debug(ev.toLine())
		return ev
	}

	xEnv.Errorf("[%s] [%s] %s %s %s %s %s %s %d %s %v",
		ev.level, ev.subject, ev.from, ev.typeof,
		ev.user, ev.auth, ev.msg, ev.rAddr, ev.rPort, ev.region, ev.err)

	return ev
	//var e string
	//if ev.err == nil {
	//	e = ""
	//} else {
	//	e = ev.err.Error()
	//}
	//xEnv.Error(ev.msg, zap.Time("time", ev.time),
	//	zap.Text("node_id", ev.id),
	//	zap.Text("inet", ev.inet),
	//	zap.Text("subject", ev.subject),
	//	zap.Text("remote_addr", ev.rAddr),
	//	zap.Int("remote_port", ev.rPort),
	//	zap.Text("region", ev.region),
	//	zap.Text("from", ev.from),
	//	zap.Text("typeof", ev.typeof),
	//	zap.Text("user", ev.user),
	//	zap.Text("auth", ev.auth),
	//	zap.Text("error", e),
	//	zap.Bool("alert", ev.alert),
	//	zap.Text("level", ev.level))
	//return ev
}

func (ev *Event) Level(i int) {
	switch i {
	case 0:
		ev.Notice()
	case 1:
		ev.Middle()
	case 2:
		ev.High()
	case 4:
		ev.Disaster()

	default:
		ev.Notice()
	}
}

func (ev *Event) Middle() *Event {
	ev.level = MIDDLE
	return ev
}

func (ev *Event) High() *Event {
	ev.level = HIGH
	return ev
}

func (ev *Event) Disaster() *Event {
	ev.level = DISASTER
	return ev
}

func (ev *Event) Notice() *Event {
	ev.level = NOTICE
	return ev
}

func (ev *Event) check() {
	max := 1024 * 1024

	if len(ev.msg) < max {
		return
	}

	if ev.err != nil {
		cch := exception.New()
		ev.E(cch.Wrap())
		return
	}
	ev.msg = ev.msg[0:max]
	ev.E(errors.New("msg data to long > 4096"))
}

func (ev *Event) CheckUpload() {
	if !ev.upload {
		return
	}

	if ev.typeof != "logger" {
		ev.upload = true
		return
	}

	switch ev.subject {

	case "发现错误":
		ev.upload = true

	case "打印信息":
		level := xEnv.LoggerLevel()
		ev.upload = level == zapcore.InfoLevel

	case "调试信息":
		level := xEnv.LoggerLevel()
		ev.upload = level == zapcore.DebugLevel
	}
}

func (ev *Event) Put() {
	ev.upload = true
	if xEnv.Adt() == nil {
		xEnv.Errorf("not found audit object")
		return
	}

	ev.check()
	ev.CheckUpload()
	xEnv.Adt().(*Audit).handle(ev)
}

func (ev *Event) Alert() *Event {
	ev.alert = true
	return ev
}

func (ev *Event) IsAlert() bool {
	return ev.alert
}

func (ev *Event) Field(key string) string {
	switch key {
	case "id":
		return ev.id
	case "inet":
		return ev.inet
	case "subject":
		return ev.subject
	case "remote_addr":
		return ev.rAddr
	case "remote_port":
		return strconv.Itoa(ev.rPort)
	case "from":
		return ev.from
	case "typeof":
		return ev.typeof
	case "user":
		return ev.user
	case "auth":
		return ev.auth
	case "msg":
		return ev.msg
	case "err":
		if ev.err == nil {
			return ""
		}
		return ev.err.Error()

	case "region":
		return ev.region

	case "alert":
		if ev.alert {
			return "true"
		} else {
			return "false"
		}
	case "up":
		if ev.upload {
			return "true"
		} else {
			return "false"
		}
	case "level":
		return ev.level

	case "raw":
		return ev.String()

	default:
		return ""
	}
}
