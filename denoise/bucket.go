package denoise

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/strutil"
)

/*
	+--------------------------------+---------------------------+-----------------------+
	|   abc_abc_abc                  |        23334              |         10            |
	+--------------------------------+---------------------------+-----------------------+
	|   abc_abc_abc                  |        23334              |         10            |
	+--------------------------------+---------------------------+-----------------------+
	|   abc_abc_abc                  |        23334              |         10            |
	+--------------------------------+---------------------------+-----------------------+
	|   abc_abc_abc                  |        23334              |         10            |
	+--------------------------------+---------------------------+-----------------------+
	|   abc_abc_abc                  |        23334              |         10            |
	+--------------------------------+---------------------------+-----------------------+

	算法:

	denoise.add{
		bucket = "aa",
		//codec  = "md5", //md5 , raw , mod
		//method = "equal",
		index  = ["cmdline" , "p_name"],
		rate   = "30/day",
	}

	denoise.add{
		bucket = "ab",
		index  = ["cmdline" , "p_name"],



	}
*/

type Bucket struct {
	co     *lua.LState
	ignore *cond.Ignore
	Entry  []*Section
}

func (bkt *Bucket) Add(s *Section) {
	bkt.Entry = append(bkt.Entry, s)
}

func (bkt *Bucket) NewElem(v interface{}) element {
	switch elem := v.(type) {
	case lua.LNilType:
		return element{eType: Noop, Value: v}
	case string:
		return element{eType: String, Raw: strutil.S2B(elem)}
	case []byte:
		return element{eType: String, Raw: elem}
	case lua.IndexEx:
		return element{eType: IndexEx, luaEx: elem}
	case cond.FieldEx:
		return element{eType: FieldEx, Field: elem}
	case lua.LValue:
		return element{eType: String, Raw: strutil.S2B(elem.String())}
	default:
		return element{eType: Noop, Value: v}
	}
}

func (bkt *Bucket) Do(v interface{}) bool {
	n := len(bkt.Entry)
	if n == 0 {
		return false
	}

	if bkt.ignore.Match(v, cond.WithCo(bkt.co)) {
		return false
	}

	elem := bkt.NewElem(v)
	if elem.eType == Noop {
		return false
	}

	for i := 0; i < n; i++ {
		if bkt.Entry[i].Do(elem) {
			return true
		}
	}
	return false
}

/*
	v := &process{
		name = "bash",
		exe  = "xxxxx",
		ssh  = "xxxxx",
	}

	v1 := v.index("name")
	v2 := v.index("exe")
	v3 := codec(v1 , v2)
	v4 := md5(v3)

	if cache.Incr(v4, ttl) > 60 {
		return true
	}

	return false
*/
