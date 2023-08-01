package env

import (
	"errors"
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"net"
	"strings"
)

var (
	regionNotFound = errors.New("not found region sdk with env")
	invalidIPAddr  = errors.New("invalid ip address")
)

func (env *Environment) withRegionL(L *lua.LState) int {
	pd := L.CheckVelaData(1)

	rv, ok := pd.Data.(vela.Region)
	if ok {
		env.WithRegion(rv)
		return 0
	}

	L.RaiseError("invalid region object")
	return 0
}

func (env *Environment) WithRegion(v interface{}) {
	if v == nil {
		env.Errorf("with region fail , got nil")
		return
	}

	rv, ok := v.(vela.Region)
	if ok {
		env.sub.region = rv
		return
	}
	env.Errorf("invalid region object type , got %T", v)
}

func (env *Environment) Region(v interface{}) (*vela.IPv4Info, error) {

	if env.sub.region == nil {
		return nil, regionNotFound
	}

	switch item := v.(type) {
	case nil:
		return nil, invalidIPAddr
	case string:
		return env.sub.region.Search(item)
	case *net.TCPAddr:
		return env.sub.region.Search(item.IP.String())
	case *net.UDPAddr:
		return env.sub.region.Search(item.IP.String())
	case net.Addr:
		return env.sub.region.Search(strings.Split(item.String(), ":")[0])

	default:
		return nil, fmt.Errorf("%v is not addr ", v)

	}
}
