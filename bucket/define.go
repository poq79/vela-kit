package bucket

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/vela"
	"go.etcd.io/bbolt"
)

type BucketEx struct {
	Num   int                          `json:"num"`
	Path  string                       `json:"path"`
	Size  int64                        `json:"size"`
	State bbolt.TxStats                `json:"state"`
	Bkt   map[string]bbolt.BucketStats `json:"bkt"`
}

func (db *Database) info(ctx *fasthttp.RequestCtx) error {
	dbname := ctx.QueryArgs().Peek("name")
	if len(dbname) == 0 {
		return fmt.Errorf("dbname got empty")
	}

	var dbx *bbolt.DB
	switch string(dbname) {
	case "ssc":
		dbx = db.ssc
	case "ssx":
		dbx = db.orm.Bolt
	default:
		return fmt.Errorf("invalid dbname")
	}

	be := &BucketEx{Bkt: make(map[string]bbolt.BucketStats)}
	err := dbx.View(func(tx *bbolt.Tx) error {
		be.Size = tx.Size()
		be.State = tx.Stats()
		be.Path = dbx.Path()

		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			be.Num++
			be.Bkt[string(name)] = b.Stats()
			return nil
		})

		return nil
	})

	if err != nil {
		return err
	}

	chunk, err := json.Marshal(be)
	if err != nil {
		return err
	}

	ctx.Write(chunk)
	return nil
}

func (db *Database) export(ctx *fasthttp.RequestCtx) error {
	args := ctx.QueryArgs()

	dbname := args.Peek("db")
	bucket := args.Peek("pkt")
	if len(dbname) == 0 || len(bucket) == 0 {
		return fmt.Errorf("not found db or bucket")
	}

	switch string(dbname) {
	case "ssc":
		bkt := &Bucket{dbx: db, export: "json", chains: [][]byte{bucket}}
		ctx.WriteString(bkt.String())
		return nil

	case "ssx":
		bkt := &Bucket{dbx: db, export: "json", chains: [][]byte{bucket}}
		ctx.WriteString(bkt.String())
		return nil
	default:
		return fmt.Errorf("not found db")
	}

}

func (db *Database) view(ctx *fasthttp.RequestCtx) error {
	err := db.ssc.View(func(tx *bbolt.Tx) error {
		_, err := tx.WriteTo(ctx)
		return err
	})

	return err
}

func (db *Database) define(r vela.Router) {
	r.GET("/api/v1/agent/bucket/info", xEnv.Then(db.info))
	r.GET("/api/v1/arr/agent/bucket/info", xEnv.Then(db.info))
	r.POST("/api/v1/arr/agent/bucket/export", xEnv.Then(db.export))
}
