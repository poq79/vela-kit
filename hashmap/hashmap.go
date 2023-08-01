package hashmap

import (
	"bytes"
	"encoding/gob"
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/mime"
)

var xEnv vela.Environment

type HMap map[string]interface{}

func Encode(v interface{}) ([]byte, error) {
	return mime.BinaryEncode(v)
}

func Decode(chunk []byte) (interface{}, error) {
	hm := newHashMap(0)
	dnc := gob.NewDecoder(bytes.NewReader(chunk))
	err := dnc.Decode(&hm)
	if err != nil {
		return nil, err
	}
	return hm, nil
}

func newHashMap(cap int) HMap {
	return make(HMap, cap)
}

/*
	struct  A {
}
*/
