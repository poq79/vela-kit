package env

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"reflect"
	"sync"
)

var (
	DefaultEncode = vela.BinaryEncode
	DefaultDecode = vela.BinaryDecode
)

type MimeHub struct {
	mux sync.RWMutex
	enc map[string]vela.EncodeFunc
	dec map[string]vela.DecodeFunc
}

func kind(v interface{}) string {

	if v == nil {
		return "nil"
	}

	vt := reflect.TypeOf(v)

LOOP:
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
		goto LOOP
	}
	return vt.String()

}

func (mhb *MimeHub) register(v interface{}, enc vela.EncodeFunc, dec vela.DecodeFunc) {
	mhb.mux.Lock()
	defer mhb.mux.Unlock()

	name := kind(v)
	if _, ok := mhb.dec[name]; ok {
		panic("duplicate mime decode name " + name)
		return
	}

	if _, ok := mhb.enc[name]; ok {
		panic("duplicate mime encode name " + name)
		return
	}
	mhb.enc[name] = enc
	mhb.dec[name] = dec
}

func (env *Environment) newMimeHub() {
	mhb := &MimeHub{
		enc: make(map[string]vela.EncodeFunc, 64),
		dec: make(map[string]vela.DecodeFunc, 64),
	}
	env.mime = mhb
}

func (env *Environment) Mime(v interface{}, enc vela.EncodeFunc, dec vela.DecodeFunc) {
	env.mime.register(v, enc, dec)
}

func (env *Environment) MimeDecode(name string, data []byte) (interface{}, error) {
	fn := env.mime.dec[name]
	if fn == nil {
		fn = DefaultDecode
	}
	return fn(data)
}

func (env *Environment) MimeEncode(v interface{}) ([]byte, string, error) {
	name := kind(v)
	fn := env.mime.enc[name]
	if fn == nil {
		fn = DefaultEncode
	}

	data, err := fn(v)
	if err == nil {
		return data, name, nil
	}
	return nil, name, err
}
