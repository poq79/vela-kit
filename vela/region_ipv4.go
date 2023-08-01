package vela

import "github.com/vela-ssoc/vela-kit/lua"

type IPv4Info struct {
	size   uint8
	cityID int64
	raw    []byte
	line   [5]uint8
	split  bool
	err    error
}

var (
	VERTICAL = []byte("|")
	byteNULL = []byte("")
)

func NewIPv4Info(cityId int64, raw []byte) *IPv4Info {
	return &IPv4Info{cityID: cityId, raw: raw, split: false}
}

func (i *IPv4Info) init() {
	if i.split {
		return
	}

	n := uint8(len(i.raw))

	var idx uint8 = 0
	var k uint8

	for k = 0; k < n; k++ {
		ch := i.raw[k]

		if ch == '|' {
			i.line[idx] = k
			idx++
		}

		if idx == 5 {
			goto done
		}
	}

done:
	if idx != 5 && i.line[idx] != n-1 {
		i.line[idx] = n
	}
	i.split = true
}

func (i *IPv4Info) Byte() []byte {
	if i == nil {
		return nil
	}

	return i.raw
}

func (i *IPv4Info) CityID() int64 {
	i.init()
	return i.cityID
}

func (i *IPv4Info) Country() []byte {
	i.init()

	if i.line[0] == 0 {
		return byteNULL
	}

	return i.raw[:i.line[0]]
}

func (i *IPv4Info) Region() []byte {
	i.init()
	if i.line[1] == 0 {
		return byteNULL
	}

	return i.raw[i.line[0]+1 : i.line[1]]
}

func (i *IPv4Info) Province() []byte {
	i.init()
	if i.line[2] == 0 {
		return byteNULL
	}

	return i.raw[i.line[1]+1 : i.line[2]]
}

func (i *IPv4Info) City() []byte {
	i.init()
	if i.line[3] == 0 {
		return byteNULL
	}

	return i.raw[i.line[2]+1 : i.line[3]]
}

func (i *IPv4Info) ISP() []byte {
	i.init()
	if i.line[4] == 0 {
		return byteNULL
	}
	return i.raw[i.line[3]+1 : i.line[4]]
}

func (i *IPv4Info) Index(L *lua.LState, key string) lua.LValue {
	if i.err != nil {
		return lua.S2L("invalid")
	}

	switch key {
	case "city":
		return lua.B2L(i.City())
	case "isp":
		return lua.B2L(i.ISP())
	case "province":
		return lua.B2L(i.Province())
	case "region":
		return lua.B2L(i.Region())
	case "country":
		return lua.B2L(i.Country())
	case "raw":
		return lua.B2L(i.raw)
	default:
		return lua.B2L(byteNULL)
	}
}
