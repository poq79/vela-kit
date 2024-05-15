package bucket

import (
	"go.etcd.io/bbolt"
	"time"
)

func NewDatabase() *Database {
	db := &Database{
		opt: &bbolt.Options{
			Timeout:        10 * time.Millisecond, //add 10ms timeout
			NoGrowSync:     true,
			NoSync:         true,
			NoFreelistSync: true,
			FreelistType:   bbolt.FreelistMapType,
		},
	}
	db.open()
	return db
}
