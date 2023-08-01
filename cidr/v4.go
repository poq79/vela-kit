package cidr

import (
	"encoding/binary"
	"gopkg.in/tomb.v2"
	"net"
)

func (t *IP) getLastIPv4() net.IP {
	if !t.v4 {
		return nil
	}

	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, t.LastIPv4())

	return ip
}

func (t *IP) FirstIPv4() uint32 {
	if !t.v4 {
		return 0
	}

	ip := t.first.To4()
	var fv uint32
	fv += uint32(ip[0]) << 24
	fv += uint32(ip[1]) << 16
	fv += uint32(ip[2]) << 8
	fv += uint32(ip[3])

	return fv
}

func (t *IP) HostCountIPv4() uint32 {
	ones, bits := t.nt.Mask.Size()
	return uint32(1 << (bits - ones))
}

func (t *IP) LastIPv4() uint32 {
	if !t.v4 {
		return 0
	}

	return t.FirstIPv4() + t.HostCountIPv4() - 1
}

func (t *IP) VisitIPv4(tom *tomb.Tomb, handle func(net.IP)) {
	end := t.LastIPv4()
	for ip := t.FirstIPv4(); ip <= end; ip++ {
		select {
		case <-tom.Dying():
			return
		default:
			handle(DtoIPv4(ip))
		}
	}
}
