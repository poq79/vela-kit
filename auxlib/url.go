package auxlib

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"net/url"
	"strconv"
	"strings"
)

// tcp://127.0.0.1:53/?port=53&exclude=53
type URL struct {
	x     *url.URL
	param url.Values
	ports []int
}

func (u *URL) Int(key string) int {
	n, _ := strconv.Atoi(u.Value(key))
	return n
}

func (u *URL) Bool(key string) bool {
	if u.Value(key) == "" {
		return false
	}

	return true
}

func (u *URL) Value(key string) string {
	if u.param == nil {
		return ""
	}

	p := u.param[key]

	if len(p) == 0 {
		return ""
	}

	return p[0]
}

var null = struct{}{}

func (u *URL) R(key string) map[int]struct{} {
	d := make(map[int]struct{})
	data := u.Value(key)
	if data == "" {
		return d
	}

	s := strings.Split(data, ",")
	for _, item := range s {
		if idx := strings.IndexByte(item, '-'); idx != -1 {
			seek, _ := strconv.Atoi(item[0:idx])
			end, _ := strconv.Atoi(item[idx+1:])

			if seek < 0 || end > 65535 || seek >= end {
				continue
			}

			for port := seek; port <= end; port++ {
				d[port] = null
			}
			continue
		}

		port, _ := strconv.Atoi(item)
		if port > 0 && port < 65535 {
			d[port] = null
		}
	}

	return d
}

func (u *URL) Port() int {
	p := u.x.Port()
	port, _ := strconv.Atoi(p)
	return port
}

func (u *URL) Ports() []int {
	ex := u.R("exclude")
	ps := u.R("port")

	r := make([]int, len(ps))
	k := 0
	for p, _ := range ps {
		if _, ok := ex[p]; ok {
			continue
		}

		r[k] = p
		k++
	}
	u.ports = r[:k]

	return u.ports
}

func (u *URL) Scheme() string {
	if u.x == nil {
		return ""
	}

	return u.x.Scheme
}

func (u *URL) Host() string {
	if u.x == nil {
		return ""
	}

	return u.x.Host
}

func (u *URL) Hostname() string {
	if u.x == nil {
		return ""
	}

	return u.x.Hostname()
}

func (u *URL) Path() string {
	if u.x == nil {
		return ""
	}
	return u.x.Path
}

func (u *URL) Request() string {
	if u.x == nil {
		return ""
	}

	return fmt.Sprintf("%s://%s%s", u.x.Scheme, u.x.Host, u.x.Path)
}

func (u *URL) parse() (err error) {
	u.param, err = url.ParseQuery(u.x.RawQuery)
	return
}

func (u *URL) String() string {
	return u.x.String()
}

func (u *URL) IsNil() bool {
	return u.x == nil
}

func (v *URL) V4() bool {
	return Ipv4(v.Hostname())
}

func (v *URL) V6() bool {
	return Ipv6(v.Hostname())
}

func CheckURL(val lua.LValue, L *lua.LState) URL {
	var u URL
	var e error

	if val.Type() != lua.LTString {
		L.RaiseError("invalid URL , got %s", val.Type().String())
		return u
	}

	u, e = NewURL(val.String())
	if e != nil {
		L.RaiseError("parse %s URL error %v", val.String(), e)
		return u
	}
	return u
}

func NewURL(raw string) (URL, error) {
	var u URL
	x, err := url.Parse(raw)
	if err != nil {
		return u, err
	}
	u = URL{x: x}
	return u, u.parse()
}
