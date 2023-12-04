package kind

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/strutil"
	"path/filepath"
	"strings"
)

func String(raw string) func(string) string {
	size := len(raw)

	return func(key string) string { // * , ext , ipv4, ipv6 , [1,3]
		switch key {
		case "*":
			return raw
		case "ext":
			return filepath.Ext(raw)
		case "ipv4":
			return strutil.String(auxlib.Ipv4(raw))
		case "ipv6":
			return strutil.String(auxlib.Ipv6(raw))
		case "ip":
			return strutil.String(auxlib.Ipv4(raw) || auxlib.Ipv6(raw))
		}

		n := len(key)
		if n < 3 {
			return raw
		}

		if key[0] != '[' {
			return raw
		}

		if key[n-1] != ']' {
			return raw
		}

		idx := strings.Index(key, ":")
		if idx < 0 {
			offset, err := auxlib.ToIntE(key[1 : n-1])
			if err != nil {
				return raw
			}

			if offset >= 1 && offset <= len(raw) {
				return string(raw[offset-1])
			}

			return raw
		}

		s := auxlib.ToInt(key[1:idx])
		e := auxlib.ToInt(key[idx+1 : n-1])
		if s > size {
			return ""
		}

		if e == 0 || e > size {
			return raw[s:]
		}

		if s > e {
			return ""
		}

		return raw[s:e]
	}
}
