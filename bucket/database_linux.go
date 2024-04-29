package bucket

import (
	"go.etcd.io/bbolt"
	"syscall"
	"time"
)

func NewDatabase() *Database {
	db := &Database{
		opt: &bbolt.Options{
			Timeout:        10 * time.Millisecond,
			NoGrowSync:     true,
			NoSync:         true,
			NoFreelistSync: true,
			FreelistType:   bbolt.FreelistMapType,
			MmapFlags:      syscall.MAP_POPULATE | syscall.MAP_NORESERVE,
		},
	}
	db.open()
	return db
}
