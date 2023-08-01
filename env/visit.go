package env

import (
	"fmt"
	tunnel "github.com/vela-ssoc/vela-tunnel"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func (env *Environment) Hide() tunnel.RawHide {
	return env.hide
}

func (env *Environment) Exe() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}

	exe, err = filepath.Abs(exe)
	return exe, err
}

// ExecDir 获取当前运行路径
func (env *Environment) ExecDir() string {
	exe, err := env.Exe()
	if err != nil {
		fmt.Printf("ssoc client got work exe fail %v", err)
		return ""
	}

	return filepath.Dir(exe)
}

func (env *Environment) Mode() string {
	return env.tab.mode
}

func (env *Environment) IsDebug() bool {
	return env.tab.mode == "debug"
}

func (env *Environment) Store(key string, v interface{}) {
	env.tupMutex.Lock()
	defer env.tupMutex.Unlock()
	if _, ok := env.tuple[key]; ok {
		//todo
	}

	env.tuple[key] = v
}

func (env *Environment) Find(key string) (interface{}, bool) {
	env.tupMutex.Lock()
	defer env.tupMutex.Unlock()
	v, ok := env.tuple[key]
	return v, ok
}

func (env *Environment) TunnelInfo() (string, string) {
	tnl, ok := env.Find("vela_tunnel_broker")
	if !ok {
		return "", ""
	}

	URL, ok := tnl.(*url.URL)
	if !ok {
		return "", ""
	}

	return URL.Hostname(), URL.Port()
}

func (env *Environment) QuietOn() {
	if env.Quiet() {
		return
	}

	env.quiet = time.Now().Unix()
}

func (env *Environment) Quiet() bool {
	if env.quiet <= 0 {
		return false
	}

	now := time.Now().Unix()

	if now-env.quiet < 3600 {
		return true
	}

	return false
}
