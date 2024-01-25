package bucket

import (
	"bytes"
	"fmt"
	auxlib "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/vela"
	"go.etcd.io/bbolt"
	"path/filepath"
)

func Tx2B(tx *Tx, name []byte, readonly bool) (*bbolt.Bucket, error) {
	if readonly {
		bkt := tx.Bucket(name)
		if bkt == nil {
			return nil, fmt.Errorf("%s not found", auxlib.B2S(name))
		}
		return bkt, nil
	}

	return tx.CreateBucketIfNotExists(name)
}

func unpack(tx *bbolt.Tx, chains [][]byte, readonly bool) (*bbolt.Bucket, error) {
	var bbt *bbolt.Bucket
	var err error

	n := len(chains)
	if n < 1 {
		return nil, fmt.Errorf("not found bucket chains")
	}

	bbt, err = Tx2B(tx, chains[0], readonly)
	if n == 1 {
		return bbt, err
	}

	//如果报错
	if err != nil {
		return bbt, err
	}

	for i := 1; i < n; i++ {
		bbt, err = Bkt2B(bbt, chains[i], readonly)
		if err != nil {
			return nil, err
		}
	}

	return bbt, nil
}

func StringsToBytes(bkt []string) [][]byte {
	n := len(bkt)
	if n == 0 {
		return nil
	}

	s := make([][]byte, n)
	for i := 0; i < n; i++ {
		s[i] = []byte(bkt[i])
	}
	return s
}

func Copy(dst, src *bbolt.DB, chains []string) error {
	bkt := StringsToBytes(chains)
	if len(bkt) == 0 {
		return nil
	}

	dstTx, err := dst.Begin(true)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			dstTx.Rollback()
		}
	}()

	// open a tx on old db for read
	tx, err := src.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	s1, err := unpack(tx, bkt, true)
	if err != nil {
		return err
	}

	d1, err := unpack(tx, bkt, false)
	if err != nil {
		return err
	}

	err = s1.ForEach(func(k, v []byte) error {
		if len(k) == 0 || len(v) == 0 {
			return nil
		}

		xEnv.Errorf("debug copy key:%s val:%s succeed", k, v)
		return d1.Put(k, v)
	})

	if err != nil {
		return err
	}

	return dstTx.Commit()
}

func Bkt2B(b *bbolt.Bucket, name []byte, readonly bool) (*bbolt.Bucket, error) {
	if readonly {
		bkt := b.Bucket(name)
		if bkt == nil {
			return nil, fmt.Errorf("%s not found", auxlib.B2S(name))
		}
		return bkt, nil
	}
	return b.CreateBucketIfNotExists(name)
}

func View(bkt *Bucket, fn func(string, interface{})) error {
	od := bkt.NewOverdue()
	err := xEnv.DB().View(func(tx *Tx) error {
		bbt, err := bkt.unpack(tx, true)
		if err != nil {
			return err
		}

		err = bbt.ForEach(func(k, v []byte) error {
			var elem element
			if e := iDecode(&elem, v); e != nil {
				xEnv.Errorf("export json %s decode error %v", k, e)
				return nil
			}
			od.IsExpire(string(k), elem)

			switch elem.mime {
			case vela.NIL:
				return nil
			case "lua.LNilType":
				return nil
			}

			iv, ie := elem.Decode()
			if ie != nil {
				xEnv.Errorf("export json %s to interface error %v", k, ie)
				return nil
			}

			fn(auxlib.B2S(k), iv)
			return nil
		})

		return err
	})

	od.clear()

	if err != nil {
		return err
	}

	return nil
}

func Bkt2Json(bkt *Bucket) []byte {
	buf := kind.NewJsonEncoder()
	buf.Tab("")
	err := View(bkt, buf.KV)
	if err != nil {
		xEnv.Errorf("export %v error %v", bkt, err)
		return nil
	}

	if err != nil {
		return nil
	}

	buf.End("}")
	return buf.Bytes()
}

func Bkt2Line(bkt *Bucket) []byte {
	var buf bytes.Buffer
	fn := func(key string, v interface{}) {
		buf.WriteString(key)
		buf.WriteByte(':')
		buf.WriteString(auxlib.ToString(v))
		buf.WriteByte('\n')
	}

	err := View(bkt, fn)
	if err != nil {
		xEnv.Errorf("export %v error %v", bkt.chains, err)
		return nil
	}

	return buf.Bytes()
}

func decode(v []byte) (interface{}, error) {
	var it element
	err := iDecode(&it, v)
	if err != nil {
		return nil, err
	}

	return it.Decode()
}

func FileTime(file string) int64 {
	base := filepath.Base(file)
	size := len(base)
	if size == 7 {
		return 0
	}

	if size != 18 {
		return -1
	}

	tv := base[5:15]
	return auxlib.ToInt64(tv)
}
