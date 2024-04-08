package bucket

import (
	"go.etcd.io/bbolt"
)

type Tx = bbolt.Tx

type Bucket struct {
	dbx      *bbolt.DB
	readOnly bool
	chains   [][]byte
	export   string
}

func (bkt *Bucket) NewOverdue() *Overdue {
	return &Overdue{db: bkt.dbx, chains: bkt.chains}
}
