package httpx

import (
	"crypto/tls"
	"net"
	"time"
)

type tlsGo struct {
	network  string
	insecure bool
	timeout  int
}

func (tg *tlsGo) Config(host string) *tls.Config {
	cfg := &tls.Config{
		InsecureSkipVerify: tg.insecure,
	}

	if host != "" {
		cfg.ServerName = host
	}

	return cfg
}

func (tg *tlsGo) dail(addr, host string) (*tls.Conn, error) {
	cfg := tg.Config(host)

	dialer := &net.Dialer{
		Timeout: time.Duration(tg.timeout) * time.Millisecond,
	}

	return tls.DialWithDialer(dialer, tg.network, addr, cfg)
}
