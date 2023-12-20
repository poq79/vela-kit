package lua

import (
	"strconv"
	"time"
)

type JsonEncoder struct {
	cache []byte
}

func Json(cap int) *JsonEncoder {
	return &JsonEncoder{cache: make([]byte, 0, cap)}
}

func (enc *JsonEncoder) Char(ch byte) {
	enc.cache = append(enc.cache, ch)
}

func (enc *JsonEncoder) Len() int {
	return len(enc.cache)
}

func (enc *JsonEncoder) append(v []byte) {
	n := len(v)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		enc.Char(v[i])
	}
}

var (
	Esc       = []byte{'\\', '\\'}
	CR        = []byte{'\\', '\r'}
	CN        = []byte{'\\', '\n'}
	CT        = []byte{'\\', '\t'}
	Quotation = []byte{'\\', '"'}
)

func (enc *JsonEncoder) WriteByte(ch byte) {
	switch ch {
	case '\\':
		enc.append(Esc)

	case '\r':
		enc.append(CR)

	case '\n':
		enc.append(CN)

	case '\t':
		enc.append(CT)
	case '"':
		enc.append(Quotation)

	default:
		enc.Char(ch)
	}
}

func (enc *JsonEncoder) WriteString(val string) {
	n := len(val)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		enc.WriteByte(val[i])
	}
}

func (enc *JsonEncoder) Write(val []byte) {
	n := len(val)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		enc.WriteByte(val[i])
	}
}

func (enc *JsonEncoder) Key(key string) {
	enc.Char('"')
	enc.WriteString(key)
	enc.Char('"')
	enc.WriteByte(':')
}

func (enc *JsonEncoder) Val(v string) {
	enc.Char('"')
	enc.WriteString(v)
	enc.Char('"')
}

func (enc *JsonEncoder) Insert(v []byte) {
	enc.Char('"')
	enc.Write(v)
	enc.Char('"')
}

func (enc *JsonEncoder) Int(n int) {
	enc.WriteString(strconv.Itoa(n))
}

func (enc *JsonEncoder) Bool(v bool) {
	if v {
		enc.Write(True)
	} else {
		enc.Write(False)
	}
}

func (enc *JsonEncoder) Long(n int64) {
	v := strconv.FormatInt(n, 10)
	enc.WriteString(v)
}

func (enc *JsonEncoder) ULong(n uint64) {
	v := strconv.FormatUint(n, 10)
	enc.WriteString(v)
}

func (enc *JsonEncoder) KT(key string, t time.Time) {
	enc.Key(key)
	enc.Val(t.String())
	enc.WriteByte(',')
}

//func (enc *JsonEncoder) KV(key , val string) {
//	enc.Key(key)
//	enc.Val(val)
//	enc.WriteByte(',')
//}

func (enc *JsonEncoder) KI(key string, n int) {
	enc.Key(key)
	enc.Int(n)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) Float64(f float64) {
	v := strconv.FormatFloat(f, 'f', -1, 64)
	enc.WriteString(v)
}

func (enc *JsonEncoder) Float32(f float32) {
	v := strconv.FormatFloat(float64(f), 'f', -1, 64)
	enc.WriteString(v)
}

func (enc *JsonEncoder) KF64(key string, v float64) {
	enc.Key(key)
	enc.Float64(v)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) KL(key string, n int64) {
	enc.Key(key)
	enc.Long(n)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) KUL(key string, n uint64) {
	enc.Key(key)
	enc.ULong(n)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) Join(key string, v []string) {
	enc.Key(key)

	enc.Arr("")
	for _, item := range v {
		enc.Val(item)
		enc.WriteByte(',')
	}

	enc.End("]")
	enc.WriteByte(',')
}

func (enc *JsonEncoder) kv1(key, v string) {
	enc.Key(key)
	enc.WriteString(v)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) kv2(key, v string) {
	enc.Key(key)
	enc.Val(v)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) KV(key string, s interface{}) {
	switch val := s.(type) {
	case nil:
		enc.kv1(key, "null")
	case bool:
		enc.kv1(key, strconv.FormatBool(val))
	case float64:
		enc.kv1(key, strconv.FormatFloat(val, 'f', -1, 64))
	case float32:
		enc.kv1(key, strconv.FormatFloat(float64(val), 'f', -1, 64))
	case int:
		enc.kv1(key, strconv.Itoa(val))
	case int64:
		enc.kv1(key, strconv.FormatInt(val, 10))
	case int32:
		enc.kv1(key, strconv.Itoa(int(val)))

	case int16:
		enc.kv1(key, strconv.FormatInt(int64(val), 10))
	case int8:
		enc.kv1(key, strconv.FormatInt(int64(val), 10))
	case uint:
		enc.kv1(key, strconv.FormatUint(uint64(val), 10))
	case uint64:
		enc.kv1(key, strconv.FormatUint(val, 10))
	case uint32:
		enc.kv1(key, strconv.FormatUint(uint64(val), 10))
	case uint16:
		enc.kv1(key, strconv.FormatUint(uint64(val), 10))
	case uint8:
		enc.kv1(key, strconv.FormatUint(uint64(val), 10))

	case string:
		enc.kv2(key, val)

	case LString:
		enc.kv2(key, string(val))
	case LBool:
		enc.kv1(key, strconv.FormatBool(bool(val)))
	case LNilType:
		enc.kv2(key, "nil")

	case LNumber:
		enc.kv1(key, strconv.FormatFloat(float64(val), 'f', -1, 64))
	case LInt:
		enc.kv1(key, strconv.Itoa(int(val)))

	case []string:
		enc.Join(key, val)
	case []byte:
		enc.kv2(key, B2S(val))

	case time.Time:
		if y := val.Year(); y < 0 || y >= 10000 {
			// RFC 3339 is clear that years are 4 digits exactly.
			// See golang.org/issue/4556#c15 for more discussion.

			return
			//return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
		}
		enc.kv2(key, val.Format(time.RFC3339Nano))
	case error:
		enc.kv2(key, val.Error())

	default:
		enc.kv2(key, "")

	}

}

var False = []byte("false")
var True = []byte("true")

func (enc *JsonEncoder) KB(key string, b bool) {
	enc.Key(key)

	if b {
		enc.Write(True)
	} else {
		enc.Write(False)
	}

	enc.WriteByte(',')
}

func (enc *JsonEncoder) False(key string) {
	enc.Key(key)
	enc.Write(False)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) True(key string) {
	enc.Key(key)
	enc.Write(True)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) Tab(name string) {
	if len(name) != 0 {
		enc.Val(name)
		enc.WriteByte(':')
	}

	enc.WriteByte('{')
}

func (enc *JsonEncoder) Arr(name string) {
	if len(name) != 0 {
		enc.Val(name)
		enc.WriteByte(':')
	}
	enc.WriteByte('[')
}

func (enc *JsonEncoder) Raw(key string, val []byte) {
	n := len(val)
	if n == 0 {
		return
	}

	if len(key) > 0 {
		enc.Key(key)
	}

	enc.append(val)
	enc.WriteByte(',')
}

func (enc *JsonEncoder) End(val string) {
	n := enc.Len()

	if n <= 0 {
		return
	}

	if enc.cache[n-1] == ',' {
		enc.cache = enc.cache[:n-1]
	}

	enc.WriteString(val)
}

func (enc *JsonEncoder) Bytes() []byte {
	return enc.cache
}

func (enc *JsonEncoder) Reset() {
	if enc.Len() > 8192*5 {
		enc.cache = nil
		return
	}
	enc.cache = enc.cache[:0]
}
