package cidr

import (
	"gopkg.in/tomb.v2"
	"math/big"
	"net"
)

func (t *IP) getLastIPv6() net.IP {
	if !t.v6 {
		return nil
	}

	ip := make(net.IP, 16)
	ip = t.LastIPv6().Bytes()
	return ip
}

// FirstIPv6 return Decimal FirstIP
func (t *IP) FirstIPv6() *big.Int {
	if !t.v6 {
		return nil
	}

	IPInt := big.NewInt(0)
	return IPInt.SetBytes(t.first.To16())
}

// HostCountIPv6 return number of IPs on the parsed range
func (t *IP) HostCountIPv6() *big.Int {
	ones, bits := t.nt.Mask.Size()
	var max = big.NewInt(1)

	return max.Lsh(max, uint(bits-ones))
}

// LastIPv6 return Decimal LastIP
func (t *IP) LastIPv6() *big.Int {
	if !t.v6 {
		return nil
	}

	IPInt := t.FirstIPv6()
	return IPInt.Add(IPInt, big.NewInt(0).Sub(t.HostCountIPv6(), big.NewInt(1)))
}

func (t *IP) VisitIPv6(tom *tomb.Tomb, handle func(net.IP)) {
	end := t.LastIPv6().Uint64()
	IPInt := new(big.Int)
	for ipv6 := t.FirstIPv6().Uint64(); ipv6 <= end; ipv6++ {
		select {
		case <-tom.Dying():
			return
		default:
			IPInt.SetUint64(ipv6)
			handle(DtoIPv6(IPInt))
		}
	}
}
