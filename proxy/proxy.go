package proxy

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/vela-ssoc/vela-kit/lua"
	"net"
)

type Proxy struct {
	addr string
}

func (p *Proxy) String() string                         { return "proxy " + p.addr }
func (p *Proxy) Type() lua.LValueType                   { return lua.LTObject }
func (p *Proxy) AssertFloat64() (float64, bool)         { return 0, false }
func (p *Proxy) AssertString() (string, bool)           { return "", false }
func (p *Proxy) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (p *Proxy) Peek() lua.LValue                       { return p }

func (p *Proxy) URL() string {
	return fmt.Sprintf("/api/v1/broker/stream/tunnel?address=%s", p.addr)
}

func (p *Proxy) Dail(ctx context.Context) (cnn net.Conn, err error) {
	conn, err := xEnv.Stream(ctx, p.URL(), nil)
	if err != nil {
		return nil, err
	}

	return &Conn{
		finished: 0,
		reader:   websocket.JoinMessages(conn, ""),
		socket:   conn,
	}, nil
}
