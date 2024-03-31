package httpx

import (
	"github.com/vela-ssoc/vela-kit/strutil"
	"net"
	"strings"
)

func WithPort(v uint16) func(r *RawHttp) {
	return func(r *RawHttp) {
		r.Port = strutil.String(v)
	}
}

func WithScheme(v string) func(r *RawHttp) {
	return func(r *RawHttp) {
		if v == "https" {
			r.TLS = true
		}
		r.Scheme = v
	}
}

func WithAddr(peer string) func(r *RawHttp) {
	return func(r *RawHttp) {
		if len(peer) == 0 {
			return
		}

		host, port, err := net.SplitHostPort(peer)
		if err != nil {
			r.Peer = peer
		}
		r.Port = port
		r.Peer = strings.TrimSpace(host)
	}

}

func WithNetConn(netConn net.Conn) func(r *RawHttp) {
	return func(r *RawHttp) {
		r.NetConn = netConn
	}
}
