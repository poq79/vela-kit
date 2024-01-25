package shared

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/reflectx"
	"time"
)

func (shm *ShareBucket) String() string                         { return "vela.share.bucket" }
func (shm *ShareBucket) Type() lua.LValueType                   { return lua.LTObject }
func (shm *ShareBucket) AssertFloat64() (float64, bool)         { return 0, false }
func (shm *ShareBucket) AssertString() (string, bool)           { return "", false }
func (shm *ShareBucket) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (shm *ShareBucket) Peek() lua.LValue                       { return shm }

func (shm *ShareBucket) set(key string, val lua.LValue, ttl int) error {
	var err error
	var expire time.Duration
	if ttl <= 0 {
		err = shm.bucket.Set(key, val)
		goto done
	}

	expire = time.Duration(ttl) * time.Millisecond
	err = shm.bucket.SetWithExpire(key, val, expire)

done:
	return err
}

func (shm *ShareBucket) LSet(L *lua.LState) int {

	key := L.CheckString(1)
	val := L.CheckAny(2)
	ttl := L.IsInt(3)
	if val.Type() == lua.LTNil {
		L.RaiseError("#2 invalid  , got nil")
		return 0
	}

	err := shm.set(key, val, ttl)
	if err != nil {
		L.Push(lua.S2L(err.Error()))
		return 1
	}
	return 0
}

func (shm *ShareBucket) LGet(L *lua.LState) int {
	key := L.CheckString(1)
	val, e := shm.bucket.Get(key)
	L.Push(reflectx.ToLValue(val, L))
	if e == nil {
		return 1
	}
	L.Push(lua.S2L(e.Error()))
	return 2
}

func (shm *ShareBucket) LDel(L *lua.LState) int {
	key := L.CheckString(1)
	shm.bucket.Remove(key)
	return 0
}

func (shm *ShareBucket) LIncr(L *lua.LState) int {
	key := L.CheckString(1)
	step := L.CheckNumber(2)
	ttl := L.IsInt(3)
	val, e := shm.bucket.Get(key)
	var ret lua.LValue
	var err error

	if e == nil {
		sum := step + ToNumber(val)
		err = shm.set(key, sum, ttl)
		ret = sum
		goto done
	}

	if e == KeyNotFoundError {
		err = shm.set(key, step, ttl)
		ret = step
		goto done
	}

	ret = lua.LNumber(0)
	err = e

done:
	L.Push(ret)
	if err == nil {
		return 1
	}
	L.Push(lua.S2L(err.Error()))
	return 2
}

func (shm *ShareBucket) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "set":
		return lua.NewFunction(shm.LSet)
	case "get":
		return lua.NewFunction(shm.LGet)
	case "del":
		return lua.NewFunction(shm.LDel)
	case "incr":
		return lua.NewFunction(shm.LIncr)
	case "count_all":
		return lua.LNumber(shm.bucket.Len(false))
	case "count":
		return lua.LNumber(shm.bucket.Len(true))
	default:
		return lua.LNil
	}
}
