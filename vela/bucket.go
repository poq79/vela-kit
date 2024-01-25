package vela

import "go.etcd.io/bbolt"

type Bucket interface {
	BatchStore(map[string]interface{}, int) error           //batch store
	Batch(func(*bbolt.Tx, *bbolt.Bucket) error, bool) error //batch api
	Store(string, interface{}, int) error                   //key , value , expire
	Replace(string, interface{}, int) error                 //key , value , expire
	Delete(string) error                                    //key
	DeleteBucket(string) error                              //bucket name
	Range(func(string, interface{}))                        //range function(key , value)
	ForEach(func(string, []byte))                           // foreach raw data
	Count() int                                             // count key
	Get(string) (interface{}, error)                        //key
	Incr(string, int, int) (int, error)
	Int(string) int
	Int64(string) int64
	Bool(string) bool
	Push(string, []byte, int64) error
	Value(string) ([]byte, error)
	Names() string
	String() string
	Clean() error
	Encode(interface{}, int) ([]byte, error) //编码对应的数据
	ReadOnly() Bucket
}
