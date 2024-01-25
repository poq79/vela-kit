package bucket

import (
	"bytes"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/reflectx"
)

var meta = map[string]*lua.LFunction{
	"byte":   lua.NewFunction(bucketMetaByte),
	"export": lua.NewFunction(bucketMetaExport),
	"get":    lua.NewFunction(BktMetaGet),
	"set":    lua.NewFunction(BktMetaSet),
	"delete": lua.NewFunction(BktMetaDel),
	"remove": lua.NewFunction(bucketMetaRemove),
	"clear":  lua.NewFunction(BktMetaClean),
	"pairs":  lua.NewFunction(bucketMetaPairs),
	"info":   lua.NewFunction(bucketMetaInfo),
	"count":  lua.NewFunction(bucketMetaCount),
	"depth":  lua.NewFunction(bucketMetaDepth),
	"incr":   lua.NewFunction(bucketMetaIncr),
	"suffix": lua.NewFunction(bucketMetaSuffix),
	"prefix": lua.NewFunction(bucketMetaPrefix),
}

func CheckBucket(L *lua.LState, idx int) *Bucket {
	obj := L.CheckObject(idx)

	bkt, ok := obj.(*Bucket)
	if ok {
		return bkt
	}
	L.RaiseError("bad argument #%d to *Bucket", idx)
	return nil
}

func BktMetaGet(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	key := L.CheckString(2)
	it, err := bkt.Load(key)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.S2L(err.Error()))
		return 1
	}

	if it.IsNil() {
		L.Push(lua.LNil)
		return 1
	}

	val, err := it.Decode()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.S2L(err.Error()))
		return 2
	}

	L.Push(reflectx.ToLValue(val, L))
	return 1
}

func BktMetaSet(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	key := L.CheckString(2)
	val := L.CheckAny(3)
	expire := L.IsInt(4)

	if e := bkt.Store(key, val, expire); e != nil {
		L.Push(lua.S2L(e.Error()))
		return 1
	}
	return 0
}

func BktMetaDel(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	n := L.GetTop()
	if n <= 1 {
		return 0
	}

	err := xEnv.DB().Batch(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, false)
		if err != nil {
			return err
		}
		errs := exception.New()
		for i := 2; i <= n; i++ {
			name := L.Get(i).String()
			errs.Try(name, bbt.Delete(lua.S2B(name)))
		}
		return errs.Wrap()
	})

	if err == nil {
		return 0
	}

	L.Push(lua.S2L(err.Error()))
	return 1
}

func BktMetaClean(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	_ = bkt.Clean()
	return 0
}

func bucketMetaRemove(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	n := L.GetTop()
	if n <= 1 {
		return 0
	}

	err := xEnv.DB().Batch(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, false)
		if err != nil {
			return err
		}
		errs := exception.New()
		for i := 2; i <= n; i++ {
			name := L.Get(i).String()
			errs.Try(name, bbt.DeleteBucket(lua.S2B(name)))
		}
		return errs.Wrap()
	})

	if err == nil {
		return 0
	}
	L.Push(lua.S2L(err.Error()))
	return 1
}

func bucketMetaPairs(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	pip := pipe.NewByLua(L, pipe.Env(xEnv), pipe.Seek(1))

	od := bkt.NewOverdue()
	err := xEnv.DB().View(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, true)
		if err != nil {
			return err
		}

		err = bbt.ForEach(func(k, v []byte) error {
			var it element
			var er error
			er = iDecode(&it, v)
			if er != nil {
				return nil
			}

			od.IsExpire(auxlib.B2S(k), it)

			if it.IsNil() {
				er = pip.Call2(lua.B2L(k), lua.LNil, L)
				return nil
			}

			lv, er := it.Decode()
			if er != nil {
				xEnv.Infof("decode bucket item error %v", er)
				return er
			}

			return pip.Call2(lua.B2L(k), lv, L)
		})

		return err
	})

	if err != nil {
		L.Pushf("%v", err)
		return 1
	}
	od.clear()

	return 0
}

func bucketMetaInfo(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	err := xEnv.DB().View(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, true)
		if err != nil {
			return err
		}

		L.Push(L.NewAnyData(bbt.Stats(), lua.Reflect(lua.ELEM)))
		return nil
	})

	if err != nil {
		xEnv.Errorf("not found %v", err)
		return 0
	}
	return 1
}

func bucketMetaCount(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	err := xEnv.DB().View(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, true)
		if err != nil {
			return err
		}
		L.Push(lua.LInt(bbt.Stats().KeyN))
		return nil
	})

	if err != nil {
		xEnv.Errorf("%v count fail", bkt)
		L.Push(lua.LInt(0))
	}

	return 1

}

func bucketMetaDepth(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	err := xEnv.DB().View(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, true)
		if err != nil {
			return err
		}

		L.Push(lua.LInt(bbt.Stats().Depth))
		return nil
	})
	if err != nil {
		xEnv.Errorf("%v count fail", bkt)
		L.Push(lua.LInt(0))
	}
	return 1
}

func bucketMetaIncr(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	key := L.CheckString(2)
	val := L.CheckNumber(3)

	expire := L.IsInt(4)
	var sum float64
	err := xEnv.DB().Batch(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, false)
		if err != nil {
			return err
		}

		b := lua.S2B(key)
		data := bbt.Get(b)
		it := &element{}
		err = iDecode(it, data)
		if err != nil {
			xEnv.Infof("incr %s decode fail error %v", key, err)
			goto INCR
		}

	INCR:
		sum = it.incr(float64(val), expire)
		return bbt.Put(b, it.Byte())
	})

	if err != nil {
		L.Push(lua.LInt(0))
		L.Push(lua.S2L(err.Error()))
		return 2
	}

	L.Push(lua.LNumber(sum))
	return 1
}

func bucketMetaExport(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	val, _ := L.Get(2).AssertString()
	bkt.export = val
	L.Push(bkt)
	return 1
}

func bucketMetaFixHelper(L *lua.LState, fn func([]byte, []byte) bool) int {
	bkt := CheckBucket(L, 1)
	fix := L.CheckString(2)
	ret := L.CreateTable(32, 0)
	od := bkt.NewOverdue()

	err := xEnv.DB().View(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, true)
		if err != nil {
			return err
		}
		i := 1

		err = bbt.ForEach(func(k, v []byte) error {
			it := element{}
			er := iDecode(&it, v)
			if er != nil {
				xEnv.Errorf("invalid item error %v", er)
				return nil
			}
			od.IsExpire(auxlib.B2S(k), it)

			if !fn(k, lua.S2B(fix)) {
				return nil
			}

			if iv, ie := it.Decode(); ie != nil {
				xEnv.Errorf("decode bucket item error %v", ie)
				goto next
			} else {
				ret.RawSetInt(i, reflectx.ToLValue(iv, L))
				return nil
			}

		next:
			ret.RawSetInt(i, lua.B2L(v))
			return nil
		})

		return err
	})

	od.clear()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.S2L(err.Error()))
		return 2
	}

	L.Push(ret)
	return 1
}
func bucketMetaSuffix(L *lua.LState) int {
	return bucketMetaFixHelper(L, bytes.HasSuffix)
}

func bucketMetaPrefix(L *lua.LState) int {
	return bucketMetaFixHelper(L, bytes.HasPrefix)
}

func bucketMetaByte(L *lua.LState) int {
	bkt := CheckBucket(L, 1)
	L.Push(lua.S2L(bkt.String()))
	return 1
}

func (bkt *Bucket) MetaTable(L *lua.LState, key string) lua.LValue {
	return meta[key]
}

func (bkt *Bucket) Index(L *lua.LState, key string) lua.LValue {
	elem, err := bkt.Load(key)
	if err != nil {
		return lua.LNil
	}

	if elem.IsNil() {
		return lua.LNil
	}

	val, err := elem.Decode()
	if err != nil {
		return lua.LNil
	}

	return reflectx.ToLValue(val, L)
}

func (bkt *Bucket) NewIndex(L *lua.LState, key string, val lua.LValue) {
	err := bkt.Store(key, val, 0)
	if err != nil {
		xEnv.Errorf("%s store %s error %v", bytes.Join(bkt.chains, []byte(",")), key, err)
	}
}

func (bkt *Bucket) Meta(L *lua.LState, key lua.LValue) lua.LValue {
	return bkt.Index(L, key.String())
}

func (bkt *Bucket) NewMeta(L *lua.LState, key lua.LValue, val lua.LValue) {
	bkt.NewIndex(L, key.String(), val)
}
