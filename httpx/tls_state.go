package httpx

import (
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
)

type state struct {
	version   uint16
	handshake bool
	IsCA      bool
	host      string
	after     int64
	subject   string
}

func (st *state) String() string                         { return lua.B2S(st.Byte()) }
func (st *state) Type() lua.LValueType                   { return lua.LTObject }
func (st *state) AssertFloat64() (float64, bool)         { return 0, false }
func (st *state) AssertString() (string, bool)           { return "", false }
func (st *state) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (st *state) Peek() lua.LValue                       { return st }

func (st state) Byte() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("version", st.version)
	enc.KV("handshake", st.handshake)
	enc.KV("is_ca", st.IsCA)
	enc.KV("host", st.host)
	enc.KV("after", st.after)
	enc.KV("subject", st.subject)
	enc.End("}")
	return enc.Bytes()
}

func (st *state) Index(L *lua.LState, key string) lua.LValue {
	switch key {

	case "version":
		return lua.LInt(st.version)

	case "handshake":
		return lua.LBool(st.handshake)

	case "is_ca":
		return lua.LBool(st.IsCA)

	case "host":
		return lua.S2L(st.host)

	case "after":
		return lua.LNumber(st.after)

	case "subject":
		return lua.S2L(st.subject)

	default:
		return lua.LNil
	}
}
