package env

import (
	"github.com/vela-ssoc/vela-kit/node"
	tunnel "github.com/vela-ssoc/vela-tunnel"
	"net"
	"runtime"
	"strconv"
)

type broker struct {
	arch   string
	mac    net.HardwareAddr
	inet   net.IP
	inet6  net.IP
	remote net.Addr
	edit   string
}

func (env *Environment) Prefix() string {
	return node.Prefix()
}

func (env *Environment) Ident() tunnel.Ident {
	if env.tnl == nil {
		return tunnel.Ident{}
	}

	return env.tnl.Ident()
}

func (env *Environment) ID() string {
	if env.tnl == nil {
		return ""
	}
	return strconv.FormatInt(env.tnl.Issue().ID, 10)
}

func (env *Environment) Arch() string {
	return runtime.GOARCH
}

func (env *Environment) Broker() (net.IP, int) {
	if env.tnl == nil {
		return nil, 0
	}

	addr, ok := env.tnl.RemoteAddr().(*net.TCPAddr)
	if ok {
		return addr.IP, addr.Port
	}
	env.Error("parse environment remote net fail")
	return nil, 0
}

func (env *Environment) Mac() string {
	if env.tnl == nil {
		return ""
	}

	return env.tnl.Ident().MAC
}

func (env *Environment) Inet() string {
	if env.tnl == nil {
		return ""
	}

	return env.tnl.Inet().String()
}

func (env *Environment) Edition() string {
	if env.tnl == nil {
		return ""
	}

	return env.tnl.Ident().Semver
}

func (env *Environment) LocalAddr() string {
	if env.tnl == nil {
		return ""
	}

	return env.tnl.LocalAddr().String()
}

func (env *Environment) RemoteAddr() string {
	if env.tnl == nil {
		return ""
	}

	return env.tnl.RemoteAddr().String()
}
