package bucket

import (
	"bytes"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/xreflect"
)

var meta = map[string]*lua.LFunction{
	"byte":   lua.NewFunction(bucketMetaByte),
	"export": lua.NewFunction(bucketMetaExport),
	"get":    lua.NewFunction(bucketMetaGet),
	"set":    lua.NewFunction(bucketMetaSet),
	"delete": lua.NewFunction(bucketMetaDelete),
	"remove": lua.NewFunction(bucketMetaRemove),
	"clear":  lua.NewFunction(bucketMetaClear),
	"pairs":  lua.NewFunction(bucketMetaPairs),
	"info":   lua.NewFunction(bucketMetaInfo),
	"count":  lua.NewFunction(bucketMetaCount),
	"depth":  lua.NewFunction(bucketMetaDepth),
	"incr":   lua.NewFunction(bucketMetaIncr),
	"suffix": lua.NewFunction(bucketMetaSuffix),
	"prefix": lua.NewFunction(bucketMetaPrefix),
}

func checkBucketValue(L *lua.LState, idx int) *Bucket {
	obj := L.CheckObject(idx)

	bkt, ok := obj.(*Bucket)
	if ok {
		return bkt
	}
	L.RaiseError("bad argument #%d to *Bucket", idx)
	return nil
}

func bucketMetaGet(L *lua.LState) int {
	bkt := checkBucketValue(L, 1)
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

	L.Push(xreflect.ToLValue(val, L))
	return 1
}

func bucketMetaSet(L *lua.LState) int {
	bkt := checkBucketValue(L, 1)
	key := L.CheckString(2)
	val := L.CheckAny(3)
	expire := L.IsInt(4)

	if e := bkt.Store(key, val, expire); e != nil {
		L.Push(lua.S2L(e.Error()))
		return 1
	}
	return 0
}

func bucketMetaDelete(L *lua.LState) int {
	bkt := checkBucketValue(L, 1)
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

func bucketMetaClear(L *lua.LState) int {
	bkt := checkBucketValue(L, 1)
	bkt.Clear()
	return 0
}

func bucketMetaRemove(L *lua.LState) int {
	bkt := checkBucketValue(L, 1)
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
	bkt := checkBucketValue(L, 1)
	pip := pipe.NewByLua(L, pipe.Env(xEnv), pipe.Seek(1))

	ee := bkt.NewExpireQueue()
	err := xEnv.DB().View(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, true)
		if err != nil {
			return err
		}

		err = bbt.ForEach(func(k, v []byte) error {
			var it item
			var er error
			er = iDecode(&it, v)
			if er != nil {
				return nil
			}

			ee.IsExpire(auxlib.B2S(k), it)

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
	ee.clear()

	return 0
}

func bucketMetaInfo(L *lua.LState) int {
	bkt := checkBucketValue(L, 1)
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
	bkt := checkBucketValue(L, 1)
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
	bkt := checkBucketValue(L, 1)
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
	bkt := checkBucketValue(L, 1)
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
		it := &item{}
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
	bkt := checkBucketValue(L, 1)
	val, _ := L.Get(2).AssertString()
	bkt.export = val
	L.Push(bkt)
	return 1
}

func bucketMetaFixHelper(L *lua.LState, fn func([]byte, []byte) bool) int {
	bkt := checkBucketValue(L, 1)
	fix := L.CheckString(2)
	ret := L.CreateTable(32, 0)
	ee := bkt.NewExpireQueue()

	err := xEnv.DB().View(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, true)
		if err != nil {
			return err
		}
		i := 1

		err = bbt.ForEach(func(k, v []byte) error {
			it := item{}
			er := iDecode(&it, v)
			if er != nil {
				xEnv.Errorf("invalid item error %v", er)
				return nil
			}
			ee.IsExpire(auxlib.B2S(k), it)

			if !fn(k, lua.S2B(fix)) {
				return nil
			}

			if iv, ie := it.Decode(); ie != nil {
				xEnv.Errorf("decode bucket item error %v", ie)
				goto next
			} else {
				ret.RawSetInt(i, xreflect.ToLValue(iv, L))
				return nil
			}

		next:
			ret.RawSetInt(i, lua.B2L(v))
			return nil
		})

		return err
	})

	ee.clear()
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
	bkt := checkBucketValue(L, 1)
	L.Push(lua.S2L(bkt.String()))
	return 1
}

func (bkt *Bucket) MetaTable(L *lua.LState, key string) lua.LValue {
	return meta[key]
}

func (bkt *Bucket) Index(L *lua.LState, key string) lua.LValue {
	it, err := bkt.Load(key)
	if err != nil {
		return lua.LNil
	}

	if it.IsNil() {
		return lua.LNil
	}

	val, err := it.Decode()
	if err != nil {
		return lua.LNil
	}

	return xreflect.ToLValue(val, L)
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
