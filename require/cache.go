package require

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"os"
	"path/filepath"
	"sync/atomic"
)

const (
	OOP CacheFlag = iota + 1
	OK
	ERR
	DEL
)

type CacheFlag uint8

type cache struct {
	status CacheFlag
	err    error
	hit    uint64
	name   string
	mtime  int64
	co     *lua.LState
	cdata  lua.LValue
}

func (c *cache) Hit() {
	atomic.AddUint64(&c.hit, 1)
}

func (c *cache) state() *lua.LState {
	if c.co == nil {
		c.co = xEnv.Coroutine()
	}

	return xEnv.Clone(c.co)
}

func (c *cache) free() {
	c.co.SetTop(0)
}

func (c *cache) file() string {
	filename := filepath.Join(xEnv.ExecDir(), "3rd", c.name)
	return filepath.Clean(filename) + ".lua"
}

func (c *cache) load() error {
	var tName = c.name + ".lua"
	info, err := xEnv.Third(tName)
	if err != nil {
		c.err = err
		c.status = ERR
		return err
	}

	return c.compile(info)
}

func (c *cache) compile(info *vela.ThirdInfo) error {
	co := xEnv.Clone(c.co)
	defer xEnv.Free(co)

	err := xEnv.DoFile(co, info.File())
	if err != nil {
		c.err = err
		c.status = ERR
		xEnv.Errorf("3rd sync compile %s error %v", info.File(), err)
		return err
	}

	c.mtime = info.MTime

	lv := co.CheckAny(-1)
	switch lv.Type() {
	case lua.LTObject, lua.LTTable, lua.LTAnyData,
		lua.LTKv, lua.LTSkv, lua.LTVelaData:
		c.cdata = lv
		c.status = OK
	default:
		err = fmt.Errorf("compile 3rd %s invalid type , must have index ,got %s", c.file(), lv.Type().String())
		c.err = err
		c.status = ERR
		xEnv.Errorf("%v", err)
	}
	return nil
}

func (c *cache) stat() int64 {
	stat, err := os.Stat(c.file())
	if err == nil {
		return stat.ModTime().Unix()
	}

	if os.IsNotExist(err) {
		return 0
	}

	xEnv.Errorf("require sync %s fail %v", c.file(), err)
	return c.mtime
}
