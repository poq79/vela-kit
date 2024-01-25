package bucket

import (
	"bytes"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/vela"
	"go.etcd.io/bbolt"
)

type Overdue struct {
	db     *bbolt.DB
	chains [][]byte
	data   []string
}

func (od *Overdue) IsExpire(key string, elem element) {
	if elem.mime != vela.EXPIRE {
		return
	}
	od.data = append(od.data, key)
}

func (od *Overdue) clear() {
	if od.db == nil {
		return
	}

	if len(od.data) == 0 {
		return
	}

	tx, err := od.db.Begin(true)
	if err != nil {
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

	bbt, err := unpack(tx, od.chains, false)
	if err != nil {
		return
	}

	for _, key := range od.data {
		if e := bbt.Delete(auxlib.S2B(key)); e != nil {
			xEnv.Errorf("%v delete %s fail %v", od.chains, key, e)
		}
	}

	err = tx.Commit()
	if err != nil {
		xEnv.Errorf("%s expire clear fail %v", bytes.Join(od.chains, []byte(".")), err)
	}
}
