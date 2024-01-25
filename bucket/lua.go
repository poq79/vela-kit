package bucket

import (
	"fmt"
	"github.com/asdine/storm/v3"
	"github.com/vela-ssoc/vela-kit/codec"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
)

var xEnv vela.Environment

func (db *Database) BucketL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.Push(lua.LNil)
		return 1
	}

	b := &Bucket{dbx: db}

	for i := 1; i <= n; i++ {
		name := L.CheckString(i)
		b.chains = append(b.chains, lua.S2B(name))
	}

	L.Push(b)
	return 1
}

func (db *Database) index(L *lua.LState, key string) lua.LValue {
	return &Bucket{
		dbx:    db,
		chains: [][]byte{[]byte(key)},
		export: "json",
	}
}

func (db *Database) infoL(L *lua.LState) int {
	tab := L.NewTable()
	i := 0
	_ = db.ssc.View(func(tx *Tx) error {
		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			i++
			tab.RawSetInt(i, lua.B2L(name))
			return nil
		})
	})
	L.Push(tab)
	return 1
}

func (db *Database) removeL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		return 0
	}

	err := db.ssc.Batch(func(tx *Tx) error {
		errs := exception.New()
		for i := 1; i <= n; i++ {
			lv := L.Get(i)
			if lv.Type() == lua.LTString {
				name := lv.String()
				errs.Try(name, tx.DeleteBucket(lua.S2B(name)))
			}
		}
		return errs.Wrap()
	})

	if err == nil {
		return 0
	}

	L.Push(lua.S2L(err.Error()))
	return 1
}

func (db *Database) compactL(L *lua.LState) int {
	db.Compact(".ssc", db.ssc, func(dst *bbolt.DB) error {
		if dst == nil {
			return fmt.Errorf("not found ssc compact db file")
		}

		/*
			Copy(dst, db.ssc, []string{"vela", "process", "snapshot"})
			Copy(dst, db.ssc, []string{"vela", "listen", "snapshot"})
			Copy(dst, db.ssc, []string{"vela", "group", "snapshot"})
			Copy(dst, db.ssc, []string{"vela", "account", "snapshot"})
			Copy(dst, db.ssc, []string{"VELA_FILE_HASH"})
			Copy(dst, db.ssc, []string{"VELA_THIRD_INFO_DB"})
			Copy(dst, db.ssc, []string{"windows_event_record_offset"})
			Copy(dst, db.ssc, []string{"windows_access_log"})
		*/

		old := db.ssc
		db.ssc = dst

		if e := old.Close(); e != nil {
			xEnv.Errorf("close %s db fail %v", old.Path(), e)
		}

		return nil
	})

	db.Compact(".ssx", db.orm.Bolt, func(dst *bbolt.DB) error {
		if dst == nil {
			return fmt.Errorf("not found ssx compact db file")
		}

		orm, err := storm.Open(dst.Path(), storm.UseDB(dst))
		if err != nil {
			return err
		}
		orm.WithCodec(codec.Sonic{})

		old := db.orm
		db.orm = orm

		if e := old.Close(); e != nil {
			xEnv.Errorf("close %s db fail %v", old.Bolt.Path(), e)
		}

		return nil
	})
	return 0
}

func (db *Database) cleanL(L *lua.LState) int {
	ms, err := filepath.Glob(fmt.Sprintf("%s/*.db", xEnv.ExecDir()))
	if err != nil {
		L.RaiseError("database glob fail %v", err)
		return 0
	}

	for _, filename := range ms {
		if filename == db.ssc.Path() || filename == db.orm.Bolt.Path() {
			continue
		}

		if e := os.Remove(filename); e != nil {
			xEnv.Errorf("%s remove fail %v", filename, e)
		} else {
			xEnv.Errorf("%s remove succeed", filename)
		}
	}
	return 0
}

func Constructor(env vela.Environment, callback func(v interface{}) error) {
	xEnv = env
	db := NewDatabase()
	if e := callback(db); e != nil {
		exception.Fatal(fmt.Errorf("callback database fail %v", e))
		return
	}

	db.define(xEnv.R())
	xEnv.Set("db",
		lua.NewExport("vela.db.export",
			lua.WithFunc(db.BucketL),
			lua.WithIndex(db.index)))

	uv := lua.NewUserKV()
	uv.Set("info", lua.NewFunction(db.infoL))
	uv.Set("remove", lua.NewFunction(db.removeL))
	uv.Set("compact", lua.NewFunction(db.compactL))
	uv.Set("clean", lua.NewFunction(db.cleanL))
	xEnv.Set("bbolt", lua.NewExport("vela.bbolt.export", lua.WithTable(uv)))
}
