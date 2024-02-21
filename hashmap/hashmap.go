package hashmap

import (
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

type HMap map[string]interface{}

func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Decode(chunk []byte) (interface{}, error) {
	hm := New(0)
	err := json.Unmarshal(chunk, &hm)
	return hm, err
}

func (hm HMap) Merge(v HMap) {
	for key, val := range v {
		hm[key] = val
	}
}

func New(cap int) HMap {
	return make(HMap, cap)
}

/*
	struct  A {
}
*/
