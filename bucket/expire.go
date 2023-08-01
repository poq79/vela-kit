package bucket

import (
	"bytes"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/vela"
	"go.etcd.io/bbolt"
)

type expireQueue struct {
	db     *bbolt.DB
	chains [][]byte
	data   []string
}

func (ee *expireQueue) IsExpire(key string, it item) {
	if it.mime != vela.EXPIRE {
		return
	}
	ee.data = append(ee.data, key)
}

func (ee *expireQueue) clear() {
	if ee.db == nil {
		return
	}

	if len(ee.data) == 0 {
		return
	}

	tx, err := ee.db.Begin(true)
	if err != nil {
		return
	}
	defer tx.Rollback()

	bbt, err := unpack(tx, ee.chains, false)
	if err != nil {
		return
	}

	for _, val := range ee.data {
		bbt.Delete(auxlib.S2B(val))
	}

	err = tx.Commit()
	if err != nil {
		xEnv.Errorf("%s expire clear fail %v", bytes.Join(ee.chains, []byte(".")), err)
	}
}
