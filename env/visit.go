package env

import (
	"fmt"
	tunnel "github.com/vela-ssoc/vela-tunnel"
	"net/url"
	"os"
	"path/filepath"
)

func (env *Environment) Hide() tunnel.RawHide {
	return env.hide
}

func (env *Environment) CPU() float64 {
	if env.rtm == nil {
		return 0
	}

	return env.rtm.CPU()
}

func (env *Environment) AgentCPU() float64 {
	if env.rtm == nil {
		return 0
	}

	return env.rtm.AgentCPU()
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

func (env *Environment) Quiet() bool {
	if env.rtm == nil {
		return false
	}

	return env.rtm.Quiet()

}
