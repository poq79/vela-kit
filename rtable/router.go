package rtable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/problem"
	"github.com/vela-ssoc/vela-kit/vela"
	tun2 "github.com/vela-ssoc/vela-tunnel"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

type TnlRouter struct {
	log    vela.Log
	mu     sync.Mutex
	hub    map[string]fasthttp.RequestHandler
	router *router.Router
	group  *router.Group
	ln     *fasthttputil.InmemoryListener
	cli    *http.Client
}

func (trr *TnlRouter) H2S() tun2.Server {
	return trr.h2s()
}

func (trr *TnlRouter) Cli() http.Client {
	return *trr.cli
}

func (trr *TnlRouter) url(req string) string {
	return fmt.Sprintf("http://ssc/%s", req)
}

func (trr *TnlRouter) Call(req string, v interface{}) (*http.Response, error) { //post
	switch data := v.(type) {
	case nil:
		return nil, nil
	case io.Reader:
		return trr.cli.Post(trr.url(req), "application/json", data)
	case string:
		reader := strings.NewReader(data)
		return trr.cli.Post(trr.url(req), "application/json", reader)
	case []byte:
		reader := bytes.NewReader(data)
		return trr.cli.Post(trr.url(req), "application/json", reader)
	case fmt.Stringer:
		reader := strings.NewReader(data.String())
		return trr.cli.Post(trr.url(req), "application/json", reader)
	case uint8, uint16, uint32, uint, uint64, int8, int16, int32, int, int64, float64, float32:
		reader := strings.NewReader(auxlib.ToString(data))
		return trr.cli.Post(trr.url(req), "application/json", reader)

	default:
		chunk, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(chunk)
		return trr.cli.Post(trr.url(req), "application/json", reader)
	}

}

func (trr *TnlRouter) Listen() {
	trr.ln = fasthttputil.NewInmemoryListener()

	go func() {
		fasthttp.Serve(trr.ln, func(ctx *fasthttp.RequestCtx) {
			trr.router.Handler(ctx)
		})
	}()

	trr.cli = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return trr.ln.Dial()
			},
		},
	}
}

func (trr *TnlRouter) Bad(ctx *fasthttp.RequestCtx, code int, opt ...func(*problem.Problem)) {
	p := problem.Problem{
		Type:     xEnv.Node(),
		Status:   code,
		Instance: string(ctx.Request.RequestURI()),
	}

	if len(opt) > 0 {
		for _, fn := range opt {
			fn(&p)
		}
	}

	body, _ := json.Marshal(p)
	ctx.Response.SetStatusCode(code)
	ctx.Write(body)
}

func (trr *TnlRouter) h2s() tun2.Server {
	return &fasthttp.Server{Handler: func(ctx *fasthttp.RequestCtx) {
		trr.router.Handler(ctx)
	}}
}

func (trr *TnlRouter) reload() {
	r := router.New()
	for key, handle := range trr.hub {
		switch {
		case strings.HasPrefix(key, fasthttp.MethodGet):
			r.GET(key[4:], handle)
			continue
		case strings.HasPrefix(key, fasthttp.MethodPost):
			r.POST(key[5:], handle)
			continue
		case strings.HasPrefix(key, fasthttp.MethodDelete):
			r.DELETE(key[7:], handle)
			continue
		case strings.HasPrefix(key, fasthttp.MethodPut):
			r.PUT(key[4:], handle)
			continue
		default:
			trr.log.Errorf("%s not allow", key)
		}
	}

	trr.router = r
}

func (trr *TnlRouter) Upsert(method string, path string, handle fasthttp.RequestHandler) error {
	trr.mu.Lock()
	defer trr.mu.Unlock()
	key := fmt.Sprintf("%s_%s", method, path)
	_, ok := trr.hub[key]
	if !ok {
		trr.hub[key] = handle
		trr.router.Handle(method, path, handle)
		return nil
	}

	if !strings.HasPrefix(path, "/api/v1/arr/lua/") {
		return fmt.Errorf("%s %s already ok", method, path)
	}

	trr.hub[key] = handle
	trr.reload()
	return nil
}

func (trr *TnlRouter) Handle(method string, path string, handle fasthttp.RequestHandler) error {
	trr.mu.Lock()
	defer trr.mu.Unlock()

	key := fmt.Sprintf("%s_%s", method, path)

	_, ok := trr.hub[key]
	if ok {
		return fmt.Errorf("%s %s already ok", method, path)
	}
	trr.hub[key] = handle
	trr.router.Handle(method, path, handle)
	return nil
}

func (trr *TnlRouter) GET(path string, handle fasthttp.RequestHandler) error {
	return trr.Handle(fasthttp.MethodGet, path, handle)
}

func (trr *TnlRouter) POST(path string, handle fasthttp.RequestHandler) error {
	return trr.Handle(fasthttp.MethodPost, path, handle)
}

func (trr *TnlRouter) DELETE(path string, handle fasthttp.RequestHandler) error {
	return trr.Handle(fasthttp.MethodDelete, path, handle)
}

func (trr *TnlRouter) Undo(method string, path string) {
	trr.mu.Lock()
	defer trr.mu.Unlock()
	key := fmt.Sprintf("%s_%s", method, path)

	_, ok := trr.hub[key]
	if !ok {
		return
	}

	delete(trr.hub, key)

	trr.reload()
}

func (trr *TnlRouter) callL(L *lua.LState) int {
	return 0
}

func (trr *TnlRouter) view(ctx *fasthttp.RequestCtx) error {
	trr.mu.Lock()
	defer trr.mu.Unlock()
	enc := kind.NewJsonEncoder()
	enc.Arr("")
	add := func(method, path string) {
		enc.Tab("")
		enc.KV("method", method)
		enc.KV("full", path)
		enc.KV("path", path[7:])
		enc.End("},")
	}
	for key, _ := range trr.hub {
		switch {
		case strings.HasPrefix(key, fasthttp.MethodGet):
			add(key[:3], key[4:])
			continue
		case strings.HasPrefix(key, fasthttp.MethodPost):
			add(key[:4], key[5:])
			continue
		case strings.HasPrefix(key, fasthttp.MethodDelete):
			add(key[:6], key[7:])
			continue
		case strings.HasPrefix(key, fasthttp.MethodPut):
			add(key[:3], key[4:])
			continue
		default:
			trr.log.Errorf("%s not allow", key)
		}
	}

	enc.End("]")
	ctx.Write(enc.Bytes())
	return nil

}
