package vela

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/problem"
	"github.com/vela-ssoc/vela-tunnel"
	"io"
	"net"
	"net/http"
)

type Router interface {
	Bad(ctx *fasthttp.RequestCtx, code int, opt ...func(problem *problem.Problem))
	Undo(string, string)
	GET(string, fasthttp.RequestHandler) error
	POST(string, fasthttp.RequestHandler) error
	DELETE(string, fasthttp.RequestHandler) error
	Handle(string, string, fasthttp.RequestHandler) error
	Cli() http.Client
	Call(url string, v interface{}) (*http.Response, error)
}

type HTTPStream interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type TnlByEnv interface { //tunnel by env
	Then(handle func(ctx *fasthttp.RequestCtx) error) func(*fasthttp.RequestCtx)
	Broker() (net.IP, int)
	R() Router
	Node() string
	Tags() []string
	Doer(prefix string) (tunnel.Doer, error)
	Oneway(path string, reader io.Reader, header http.Header) error
	Fetch(path string, reader io.Reader, header http.Header) (*http.Response, error)
	JSON(path string, data interface{}, result interface{}) error
	Push(path string, data interface{}) error
	OnConnect(name string, todo func() error)
	Stream(context.Context, string, http.Header) (*websocket.Conn, error)
	Attachment(name string) (*tunnel.Attachment, error)
}
