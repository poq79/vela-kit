package cidr

import (
	"encoding/binary"
	"gopkg.in/tomb.v2"
	"math/big"
	"net"
)

type IP struct {
	first net.IP
	last  net.IP
	nt    *net.IPNet
	v4    bool
	v6    bool
}

func Parse(raw string) (*IP, error) {
	fv, nt, er := net.ParseCIDR(raw)
	if er != nil {
		return nil, er
	}

	var v4, v6 bool

	if fv.To4() != nil {
		v4 = true
	} else {
		v6 = true
	}

	tv := IP{
		nt:    nt,
		v4:    v4,
		v6:    v6,
		first: fv,
	}

	if tv.v4 {
		tv.last = tv.getLastIPv4()
	} else {
		tv.last = tv.getLastIPv6()
	}

	return &tv, nil
}

func IPv4tod(ip net.IP) uint32 {
	if ip.To4() == nil {
		return 0
	}

	ip = ip.To4()
	var intIP uint32
	intIP += uint32(ip[0]) << 24
	intIP += uint32(ip[1]) << 16
	intIP += uint32(ip[2]) << 8
	intIP += uint32(ip[3])

	return intIP
}

func IPv6tod(ip net.IP) *big.Int {
	if ip.To4() != nil {
		return nil
	}
	return big.NewInt(0).SetBytes(ip.To16())
}

func DtoIPv4(i uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, i)

	return ip
}

func DtoIPv6(i *big.Int) net.IP {
	ip := make(net.IP, 16)
	ip = i.Bytes()

	return ip
}

func Visit(tom *tomb.Tomb, tt []*IP, handle func(net.IP)) {
	if handle == nil {
		return
	}

	n := len(tt)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		select {
		case <-tom.Dying():
			return

		default:
			t := tt[i]
			if t.v4 {
				t.VisitIPv4(tom, handle)
			} else {
				t.VisitIPv6(tom, handle)
			}
		}
	}
}
