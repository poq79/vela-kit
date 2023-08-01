package bucket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/mime"
	"github.com/vela-ssoc/vela-kit/vela"
	"strconv"
	"time"
)

type item struct {
	size  uint64
	ttl   uint64
	mime  string
	chunk []byte
}

func (it *item) set(name string, chunk []byte, expire int) {
	var ttl uint64

	if expire > 0 {
		ttl = uint64(time.Now().UnixMilli()) + uint64(expire)
	} else {
		ttl = 0
	}

	//如果ttl 为空 第二次传值有时间
	if it.ttl == 0 {
		it.ttl = ttl
	}

	it.mime = name
	it.size = uint64(len(name))
	it.chunk = chunk

}

func iEncode(it *item, v interface{}, expire int) error {
	chunk, name, err := mime.Encode(v)
	if err != nil {
		return err
	}
	it.set(name, chunk, expire)
	return nil
}

func iDecode(it *item, data []byte) error {
	n := len(data)
	if n == 0 {
		it.mime = vela.NIL
		it.size = 3
		it.chunk = nil
		return nil
	}

	if n < 16 {
		return fmt.Errorf("inavlid item , too small")
	}

	size := binary.BigEndian.Uint64(data[0:8])
	ttl := binary.BigEndian.Uint64(data[8:16])
	now := time.Now().UnixMilli()

	if ttl == 0 || int64(ttl) > now {
		if size+16 == uint64(n) {
			return fmt.Errorf("inavlid item , too big")
		}

		name := data[16 : 16+size]
		chunk := data[size+16:]

		it.size = size
		it.ttl = ttl
		it.mime = auxlib.B2S(name)
		it.chunk = chunk
		return nil
	}

	it.size = 3
	it.mime = vela.EXPIRE
	it.chunk = it.chunk[:0]
	it.ttl = 0
	return nil
}

func (it item) Byte() []byte {
	var buf bytes.Buffer
	size := make([]byte, 8)
	binary.BigEndian.PutUint64(size, it.size)
	buf.Write(size)

	ttl := make([]byte, 8)
	binary.BigEndian.PutUint64(ttl, it.ttl)
	buf.Write(ttl)

	buf.WriteString(it.mime)
	buf.Write(it.chunk)
	return buf.Bytes()
}

func (it item) Decode() (interface{}, error) {
	if it.mime == "" {
		return nil, fmt.Errorf("not found mime type")
	}

	if it.IsNil() {
		return nil, nil
	}

	return mime.Decode(it.mime, it.chunk)
}

func (it item) IsNil() bool {
	return it.size == 0 || it.mime == vela.NIL || it.mime == vela.EXPIRE
}

func (it *item) incr(v float64, expire int) (sum float64) {
	num, err := it.Decode()
	if err != nil {
		xEnv.Errorf("mime: %s chunk: %s decode fail", it.mime, it.chunk)
		goto NEW
	}

	switch n := num.(type) {
	case nil:
		sum = v
	case float64:
		sum = n + v
	case float32:
		sum = float64(n) + v
	case int:
		sum = float64(n) + v
	case int8:
		sum = float64(n) + v
	case int16:
		sum = float64(n) + v
	case int32:
		sum = float64(n) + v
	case int64:
		sum = float64(n) + v
	case uint:
		sum = float64(n) + v
	case uint8:
		sum = float64(n) + v
	case uint16:
		sum = float64(n) + v
	case uint32:
		sum = float64(n) + v
	case uint64:
		sum = float64(n) + v
	case lua.LNumber:
		sum = float64(n) + v
	case lua.LInt:
		sum = float64(n) + v
	case string:
		nf, _ := strconv.ParseFloat(n, 10)
		sum = nf + v
	case []byte:
		nf, _ := strconv.ParseFloat(auxlib.B2S(n), 10)
		sum = nf + v

	default:
		sum = v
	}

NEW:
	chunk, name, _ := mime.Encode(sum)
	it.set(name, chunk, expire)
	return
}
