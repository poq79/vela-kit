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

type element struct {
	size  uint64
	ttl   uint64
	mime  string
	chunk []byte
}

func (elem *element) set(name string, chunk []byte, expire int) {
	var ttl uint64

	if expire > 0 {
		ttl = uint64(time.Now().UnixMilli()) + uint64(expire)
	} else {
		ttl = 0
	}

	//如果ttl 为空 第二次传值有时间
	if elem.ttl == 0 {
		elem.ttl = ttl
	}

	elem.mime = name
	elem.size = uint64(len(name))
	elem.chunk = chunk
}

func iEncode(elem *element, v interface{}, expire int) error {
	chunk, name, err := mime.Encode(v)
	if err != nil {
		return err
	}
	elem.set(name, chunk, expire)
	return nil
}

func iDecode(elem *element, data []byte) error {
	n := len(data)
	if n == 0 {
		elem.mime = vela.NIL
		elem.size = 3
		elem.chunk = nil
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

		elem.size = size
		elem.ttl = ttl
		elem.mime = auxlib.B2S(name)
		elem.chunk = chunk
		return nil
	}

	elem.size = 3
	elem.mime = vela.EXPIRE
	elem.chunk = elem.chunk[:0]
	elem.ttl = 0
	return nil
}

func (elem element) Byte() []byte {
	var buf bytes.Buffer
	size := make([]byte, 8)
	binary.BigEndian.PutUint64(size, elem.size)
	buf.Write(size)

	ttl := make([]byte, 8)
	binary.BigEndian.PutUint64(ttl, elem.ttl)
	buf.Write(ttl)

	buf.WriteString(elem.mime)
	buf.Write(elem.chunk)
	return buf.Bytes()
}

func (elem element) Decode() (interface{}, error) {
	if elem.mime == "" {
		return nil, fmt.Errorf("not found mime type")
	}

	if elem.IsNil() {
		return nil, nil
	}

	return mime.Decode(elem.mime, elem.chunk)
}

func (elem element) IsNil() bool {
	return elem.size == 0 || elem.mime == vela.NIL || elem.mime == vela.EXPIRE
}

func (elem *element) incr(v float64, expire int) (sum float64) {
	num, err := elem.Decode()
	if err != nil {
		xEnv.Errorf("mime: %s chunk: %s decode fail", elem.mime, elem.chunk)
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
	elem.set(name, chunk, expire)
	return
}
