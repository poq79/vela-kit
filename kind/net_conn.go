package kind

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"net"
)

type Conn struct {
	Raw  []byte
	conn net.Conn
}

func (cn Conn) String() string                         { return lua.B2S(cn.Byte()) }
func (cn Conn) Type() lua.LValueType                   { return lua.LTObject }
func (cn Conn) AssertFloat64() (float64, bool)         { return 0, false }
func (cn Conn) AssertString() (string, bool)           { return "", false }
func (cn Conn) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (cn Conn) Peek() lua.LValue                       { return cn }

func (cn Conn) Byte() []byte {
	enc := NewJsonEncoder()
	enc.Tab("")
	dst, port := cn.Destination()
	enc.KV("dst", dst)
	enc.KV("dst_port", port)

	src, port := cn.Source()
	enc.KV("src", src)
	enc.KV("src_port", port)
	enc.KV("destination", cn.conn.RemoteAddr().String())
	enc.KV("source", cn.conn.LocalAddr().String())
	enc.End("}")
	return enc.Bytes()
}

func AddrToVisual(addr net.Addr) (ip string, port int) {
	switch v := addr.(type) {
	case *net.TCPAddr:
		port = v.Port
		ip = v.IP.String()
	case *net.UDPAddr:
		port = v.Port
		ip = v.IP.String()
	}
	return

}

func (cn Conn) Source() (ip string, port int) {
	ip, port = AddrToVisual(cn.conn.LocalAddr())
	return
}

func (cn Conn) Destination() (ip string, port int) {
	ip, port = AddrToVisual(cn.conn.RemoteAddr())
	return
}

func (cn Conn) closeL(L *lua.LState) int {
	e := cn.conn.Close()
	if e != nil {
		L.Push(lua.S2L(e.Error()))
		return 1
	}
	return 0
}

func (cn Conn) Index(L *lua.LState, key string) lua.LValue {
	switch key {

	case "push":
		return L.NewFunction(cn.Push)
	case "pushf":
		return L.NewFunction(cn.Pushf)

	case "dst":
		ip, _ := cn.Destination()
		return lua.S2L(ip)

	case "src":
		ip, _ := cn.Source()
		return lua.S2L(ip)

	case "src_port":
		_, port := cn.Source()
		return lua.LInt(port)

	case "dst_port":
		_, port := cn.Destination()
		return lua.LInt(port)

	case "source":
		return lua.S2L(cn.conn.LocalAddr().String())

	case "destination":
		return lua.S2L(cn.conn.RemoteAddr().String())

	case "close":
		return lua.NewFunction(cn.closeL)

	default:
		return lua.LNil
	}
}

func (cn Conn) Push(L *lua.LState) int {
	if cn.conn == nil {
		return 0
	}

	n := L.GetTop()
	if n == 0 {
		return 0
	}

	for i := 1; i <= n; i++ {
		val := L.Get(i)
		if val.Type() == lua.LTNil {
			return 0
		}

		_, e := cn.conn.Write(lua.S2B(L.Get(i).String()))
		if e != nil {
			L.Push(lua.S2L(e.Error()))
			return 1
		}
	}

	return 0
}

func (cn Conn) Pushf(L *lua.LState) int {
	if cn.conn == nil {
		return 0
	}

	n := L.GetTop()
	if n == 0 {
		return 0
	}

	chunk := auxlib.Format(L, 0)
	_, e := cn.conn.Write(lua.S2B(chunk))
	if e != nil {
		return 0
	}
	L.Push(lua.S2L(e.Error()))
	return 1
}

func NewConn(conn net.Conn) Conn {
	return Conn{
		conn: conn,
	}
}
