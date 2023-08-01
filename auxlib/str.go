package auxlib

import (
	"errors"
	"net"
	"reflect"
	"unsafe"
)

var (
	invalidName   = errors.New("invalid name")
	invalidNameCh = errors.New("name start must char")

	invalidWrap = errors.New("invalid warp")
)

func S2B(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func IsInt(ch byte) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

func IsChar(ch byte) bool {
	if ch >= 'a' && ch <= 'z' {
		return true
	}

	if ch >= 'A' && ch <= 'Z' {
		return true
	}

	return false
}

func Name(v string) error {
	if len(v) < 2 {
		return invalidName
	}

	if !IsChar(v[0]) {
		return invalidName
	}

	n := len(v)
	for i := 1; i < n; i++ {
		ch := v[i]
		switch {
		case IsChar(ch), IsInt(ch):
			continue

		case ch == '_':
			continue

		default:
			return invalidName
		}
	}

	return nil
}

func Warp(v string) error {
	if v != "\n" && v != "\r\n" {
		return invalidWrap
	}
	return nil
}

func Ipv4(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}

	for i := 0; i < len(addr); i++ {
		if addr[i] == '.' {
			return true
		}
	}

	return false
}

func Ipv6(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}

	for i := 0; i < len(addr); i++ {
		if addr[i] == ':' {
			return true
		}
	}

	return false
}
