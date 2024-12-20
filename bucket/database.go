package bucket

import (
	"fmt"
	"github.com/asdine/storm/v3"
	"github.com/vela-ssoc/vela-kit/codec"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"go.etcd.io/bbolt"
	"path/filepath"
	"sync/atomic"
	"time"
)

type Database struct {
	compacting uint32
	opt        *bbolt.Options
	ssc        *bbolt.DB
	orm        *storm.DB
	shm        *bbolt.DB
}

func (db *Database) walk(name string) string {
	dir := xEnv.ExecDir()
	pattern := dir + fmt.Sprintf("/%s*.db", name)

	ms, err := filepath.Glob(pattern)
	if err != nil {
		return filepath.Join(dir, fmt.Sprintf("%s.db", name))
	}

	n := len(ms)
	if n == 0 {
		return filepath.Join(dir, fmt.Sprintf("%s.db", name))
	}

	var max int64
	var last string
	for i := 0; i < n; i++ {
		file := ms[i]
		ft := FileTime(file)
		if ft >= max {
			max = ft
			last = file
		}
	}

	if last != "" {
		return last
	}

	return filepath.Join(dir, fmt.Sprintf("%s.db", name))
}

func (db *Database) filepath(name string) string {
	now := time.Now().Unix()
	return filepath.Join(xEnv.ExecDir(), fmt.Sprintf("%s-%d.db", name, now))
}

func (db *Database) open() {

	//新建数据存储
	path := db.walk(".ssc")
	ssc, err := bbolt.Open(path, 0600, db.opt)
	exception.Fatal(err)
	db.ssc = ssc

	path = db.walk(".ssx")
	orm, err := storm.Open(path, storm.BoltOptions(0600, db.opt))
	exception.Fatal(err)
	orm.WithCodec(codec.Sonic{})
	db.orm = orm

	//新建数据存储
	path = db.walk(".shm")
	shm, err := bbolt.Open(path, 0600, db.opt)
	exception.Fatal(err)
	db.shm = shm
}

func (db *Database) DB() *bbolt.DB {
	return db.ssc
}

func (db *Database) ShmDB() *bbolt.DB { return db.shm }

func (db *Database) Close() error {
	errs := exception.New()
	errs.Try(".ssc", db.ssc.Close())
	errs.Try(".ssx", db.ssc.Close())
	errs.Try(".shm", db.ssc.Close())
	return errs.Wrap()
}

func (db *Database) Bucket(v ...string) vela.Bucket {
	n := len(v)
	if n == 0 {
		return nil
	}

	b := &Bucket{dbx: db.ssc, export: "json"}

	for i := 0; i < n; i++ {
		b.chains = append(b.chains, lua.S2B(v[i]))
	}
	return b
}

func (db *Database) Shm(v ...string) vela.Bucket {
	n := len(v)
	if n == 0 {
		return nil
	}
	b := &Bucket{dbx: db.shm, export: "json"}
	for i := 0; i < n; i++ {
		b.chains = append(b.chains, lua.S2B(v[i]))
	}
	return b
}

func (db *Database) Storm(v ...string) storm.Node {
	return db.orm.From(v...)
}

func (db *Database) Compacting() bool {
	return atomic.AddUint32(&db.compacting, 1) > 1
}

func (db *Database) Compact(name string, src *bbolt.DB, callback func(*bbolt.DB) error) {
	if db.Compacting() {
		xEnv.Errorf("%s database compacting", name)
		return
	}
	defer atomic.StoreUint32(&db.compacting, 0)

	path := db.filepath(name)
	dbx, err := bbolt.Open(path, 0600, db.opt)
	if err != nil {
		xEnv.Errorf("%s database open fail %v", path, err)
		return
	}

	err = bbolt.Compact(dbx, src, 0)
	if err != nil {
		xEnv.Errorf("%s compact fail %v", path, err)
		return
	}

	err = callback(dbx)
	if err != nil {
		xEnv.Errorf("%s compact callback fail %v", name, err)
	}
}
