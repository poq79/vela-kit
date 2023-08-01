package env

import (
	"github.com/asdine/storm/v3"
	"github.com/vela-ssoc/vela-kit/vela"
	"go.etcd.io/bbolt"
)

type database interface {
	DB() *bbolt.DB
	Bucket(...string) vela.Bucket
	Storm(...string) storm.Node
}

func (env *Environment) Bucket(v ...string) vela.Bucket {
	return env.db.Bucket(v...)
}

func (env *Environment) Storm(v ...string) storm.Node {
	return env.db.Storm(v...)
}

func (env *Environment) DB() *bbolt.DB {
	return env.db.DB()
}
