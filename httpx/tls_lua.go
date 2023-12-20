package httpx

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"strings"
)

func (tg *tlsGo) String() string                         { return fmt.Sprintf("web.tls.client %p", tg) }
func (tg *tlsGo) Type() lua.LValueType                   { return lua.LTObject }
func (tg *tlsGo) AssertFloat64() (float64, bool)         { return 0, false }
func (tg *tlsGo) AssertString() (string, bool)           { return "", false }
func (tg *tlsGo) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (tg *tlsGo) Peek() lua.LValue                       { return tg }

func (tg *tlsGo) e(L *lua.LState, err error) int {
	L.Push(lua.LNil)
	L.Push(lua.S2L(err.Error()))
	return 2
}

func (tg *tlsGo) dailL(L *lua.LState) int {
	n := L.GetTop()
	var addr string
	var host string
	switch n {
	case 0:
		L.RaiseError("tls client dail invalid , must be dail(addr , [hostname])")
		return 0
	case 1:
		addr = L.CheckString(1)
		host = strings.Split(addr, ":")[0]
	case 2:
		addr = L.CheckString(1)
		host = L.CheckString(2)
	default:
		L.RaiseError("tls client dail too many %d , must be dail(addr , [hostname])", n)
		return 0
	}

	conn, err := tg.dail(addr, host)
	if err != nil {
		return tg.e(L, err)
	}
	defer conn.Close()

	st := conn.ConnectionState()
	peer := st.PeerCertificates[0]
	sv := state{
		version: st.Version,
		IsCA:    peer.IsCA,
		host:    st.ServerName,
		after:   peer.NotAfter.Unix(),
		subject: peer.Subject.String(),
	}

	L.Push(&sv)
	return 1

}

func (tg *tlsGo) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "dail":
		return lua.NewFunction(tg.dailL)

	}

	return lua.LNil
}

func (tg *tlsGo) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "network":
		tg.network = val.String()
	case "timeout":
		tg.timeout = lua.CheckInt(L, val)
	case "insecure":
		tg.insecure = lua.CheckBool(L, val)
	}
}

func newLuaTlsInfo(L *lua.LState) int {

	tg := &tlsGo{network: "tcp", timeout: 1000, insecure: true}

	n := L.GetTop()
	if n == 0 {
		L.Push(tg)
		return 1
	}

	tab := L.CheckTable(1)
	tab.Range(func(key string, val lua.LValue) {
		tg.NewIndex(L, key, val)
	})
	L.Push(tg)
	return 1
}
