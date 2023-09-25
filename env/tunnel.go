package env

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/problem"
	"github.com/vela-ssoc/vela-kit/safecall"
	"github.com/vela-ssoc/vela-kit/vela"
	tunnel "github.com/vela-ssoc/vela-tunnel"
	"io"
	"net/http"
	"os"
	"time"
)

var notFoundTnlE = fmt.Errorf("not found tunnel")

type TunnelFunc func(*Environment, *fasthttp.RequestCtx) error

func after(env *Environment) {
	size := len(env.onConnect)
	if size == 0 {
		return
	}

	for i := 0; i < size; i++ {
		item := env.onConnect[i]
		env.Errorf("%s onconnect todo start", item.name)
		func(name string, todo func() error) {
			safecall.New(true).
				Timeout(60 * time.Second).
				OnError(func(err error) { env.Errorf("%s on connect todo exec fail %v", name, err) }).
				OnTimeout(func() { env.Errorf("%s on connect todo exec timeout", name) }).
				OnPanic(func(v interface{}) { env.Errorf("%s on connect todo exec panic %v", name, v) }).
				Exec(todo)
			env.Errorf("%s onconnect todo exec over", item.name)
		}(item.name, item.todo)
	}
}

func (env *Environment) Oneway(path string, reader io.Reader, header http.Header) error {
	if env.tnl == nil {
		return notFoundTnlE
	}
	return env.tnl.Oneway(env.Context(), path, reader, header)
}

func (env *Environment) Fetch(path string, reader io.Reader, header http.Header) (*http.Response, error) {
	if env.tnl == nil {
		return nil, notFoundTnlE
	}
	return env.tnl.Fetch(env.Context(), path, reader, header)
}

func (env *Environment) JSON(path string, data interface{}, result interface{}) error {
	if env.tnl == nil {
		return notFoundTnlE
	}
	return env.tnl.JSON(env.Context(), path, data, result)
}

func (env *Environment) Push(path string, data interface{}) error {
	if env.tnl == nil {
		return notFoundTnlE
	}

	return env.tnl.OnewayJSON(env.Context(), path, data)
}

func (env *Environment) Stream(ctx context.Context, path string, header http.Header) (*websocket.Conn, error) {
	if env.tnl == nil {
		return nil, notFoundTnlE
	}

	return env.tnl.Stream(ctx, path, header)
}

func (env *Environment) Attachment(addr string) (*tunnel.Attachment, error) {
	if env.tnl == nil {
		return nil, notFoundTnlE
	}
	return env.tnl.Attachment(env.Context(), addr)
}

func (env *Environment) Tags() []string {
	return env.hide.Tags
}

func (env *Environment) OnConnect(name string, todo func() error) {
	for _, ev := range env.onConnect {
		if ev.name == name {
			env.Errorf("%s On tunnel connect function already ok", name)
			return
		}
	}

	env.onConnect = append(env.onConnect, onConnectEv{
		name: name,
		todo: todo,
	})

}

func (env *Environment) Disconnect(e error) {}

func (env *Environment) Reconnected(addr *tunnel.Address) {
}

func (env *Environment) Shutdown(err error) {
}

func (env *Environment) Then(fn func(ctx *fasthttp.RequestCtx) error) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		env.Infof("[INTO] ----------------: %s", ctx.Request.URI())
		err := fn(ctx)
		env.Infof("[OVER] ----------------: %s", ctx.Request.URI())
		if err != nil {
			env.router.Bad(ctx, http.StatusInternalServerError, problem.Title("内部错误"), problem.Detail("%v", err))
			return
		}
	}
}

func (env *Environment) Worker() {
	// 从自身文件中取出携带的配置, 测试环境, 不作错误处理
	exe, err := env.Exe()
	if err != nil {
		env.Errorf("not found executable file %v", err)
		return
	}

	raw, hide, err := tunnel.ReadHide(exe)
	if err != nil {
		env.Errorf("executable read binary hide info fail %v", err)
		env.Kill(os.Kill)
		os.Exit(0)
		return
	}

	env.Debugf("read hide raw %s", raw)
	tnl, err := tunnel.Dial(env.Context(), hide, env.router.H2S(), tunnel.WithLogger(env), tunnel.WithInterval(time.Second*30), tunnel.WithNotifier(env))

	if err != nil {
		env.Errorf("vela minion tunnel v2 connect fail %v", err)
		return
	}

	env.tnl = tnl
	env.hide = raw
	env.startup()
	after(env)
}

func (env *Environment) Dev(lan string, vip string, edit string, host string) {

	tnl, err := tunnel.Dial(env.Context(), tunnel.Hide{
		Semver:   edit,
		Ethernet: tunnel.Addresses{{Addr: lan, Name: host}},
		Internet: tunnel.Addresses{{Addr: vip, Name: host}},
	}, env.router.H2S(), tunnel.WithLogger(env), tunnel.WithInterval(time.Second*30), tunnel.WithNotifier(env))

	if err != nil {
		env.Errorf("vela minion tunnel v2 connect fail %v", err)
		return
	}

	env.tnl = tnl
	env.startup()
	after(env)
}

func (env *Environment) R() vela.Router {
	return env.router
}

func (env *Environment) Node() string {
	return env.tnl.NodeName()
}

func (env *Environment) Doer(prefix string) (tunnel.Doer, error) {
	if env.tnl == nil {
		return nil, notFoundTnlE
	}

	doer := env.tnl.Doer(prefix)
	if doer == nil {
		return nil, fmt.Errorf("tunnel doer %s fail", prefix)
	}

	return doer, nil
}
