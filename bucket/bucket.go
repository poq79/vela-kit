package bucket

import (
	"go.etcd.io/bbolt"
)

type Tx = bbolt.Tx

type Bucket struct {
	dbx      *Database
	readOnly bool
	chains   [][]byte
	export   string
}

func (bkt *Bucket) NewExpireQueue() *expireQueue {
	return &expireQueue{db: bkt.dbx.ssc, chains: bkt.chains}
}
