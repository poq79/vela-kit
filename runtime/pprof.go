package runtime

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/worker"
	"github.com/vela-ssoc/vela-kit/lua"
	"gopkg.in/tomb.v2"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func pprofL(L *lua.LState) int {
	url := auxlib.CheckURL(L.Get(1), L)
	if url.IsNil() {
		L.RaiseError("bind fail")
		return 0
	}

	tom := new(tomb.Tomb)
	task := func() {
		ln, err := net.Listen(url.Scheme(), url.Host())
		if err != nil {
			xEnv.Errorf("pprof listen fail %v", err)
			return
		}
		defer ln.Close()

		go http.Serve(ln, nil)
		<-tom.Dying()
	}

	kill := func() {
		tom.Kill(fmt.Errorf("over"))
	}

	w := worker.New(L, "pprof").Env(xEnv).Task(task).Kill(kill)
	xEnv.Start(L, w).From(L.CodeVM()).Do()
	return 0
}
