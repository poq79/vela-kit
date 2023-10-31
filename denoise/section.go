package denoise

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/vela-ssoc/vela-kit/cacheutil"
	"github.com/vela-ssoc/vela-kit/strutil"
	"time"
)

var separator = []byte("_")

type Section struct {
	Index []string
	Count int
	TTL   int
	Cache *cacheutil.Cache
}

func NewSection(size int) *Section {
	return &Section{
		TTL:   5 * 60, //5分钟
		Cache: cacheutil.NewCache(size),
	}

}

func (sec *Section) key(elem element) []byte {
	n := len(sec.Index)
	if n == 0 || elem.eType == String {
		return elem.Raw
	}

	hash := md5.New()
	noop := 0
	for i := 0; i < n; i++ {
		val := elem.v(sec.Index[i])
		if len(val) == 0 {
			noop++
		}
		if i > 0 {
			hash.Write(separator)
		}
		hash.Write(val)
	}

	if noop == n {
		return nil
	}

	key := hex.EncodeToString(hash.Sum(nil))
	return strutil.S2B(key)
}

func (sec *Section) Incr(key []byte, tv int) int {
	ret := 0
	sec.Cache.Update(key, func(value []byte, found bool, expireAt uint32) ([]byte, bool, int) {
		if !found {
			ret = 1
			return strutil.Uint64(uint64(1)), true, tv
		}

		now := uint32(time.Now().Unix())
		ttl := expireAt - now
		if ttl < 0 {
			ret = 1
			return strutil.Uint64(uint64(1)), true, tv
		}
		val := strutil.B2U32(value) + uint32(1)

		ret = int(val)
		return strutil.Uint32(val), true, int(ttl)
	})
	return ret
}

func (sec *Section) Do(elem element) bool {
	key := sec.key(elem)
	if key == nil {
		return false
	}

	cnt := sec.Incr(key, sec.TTL)

	fmt.Printf("key: %s  cnt:%d\n", strutil.B2S(key), cnt)
	if cnt > sec.Count {
		return true
	}

	return false
}
