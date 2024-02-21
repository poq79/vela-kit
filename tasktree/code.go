package tasktree

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/radix"
	"github.com/vela-ssoc/vela-kit/vela"
	"strings"
	"sync/atomic"

	//"sync"
	"time"
)

type header struct {
	id int64

	key string

	hash string

	uptime time.Time

	mask []string

	link []string

	way vela.Way

	disable bool //不允许创建

	status uint32

	env vela.Environment

	dialect bool

	err error
}

type Code struct {
	//pcall 状态
	cas uint32

	//头信息
	header *header

	//tree 关键信息
	tree radix.Tree

	cancel context.CancelFunc

	//fn     *lua.LFunction
	chunk []byte

	//meta
	metadata map[string]interface{}

	//param
	param map[string]interface{}
}

func newCodeVM(key string, chunk []byte, env vela.Environment, way vela.Way) (*lua.LState, *Code) {
	ctx, cancel := context.WithCancel(xEnv.Context())
	code := &Code{
		header: &header{
			hash:   fmt.Sprintf("%x", md5.Sum(chunk)),
			status: uint32(vela.Register),
			way:    way,
			key:    key,
			env:    env,
		},
		tree:   radix.NewSafeTree(),
		cancel: cancel,
		chunk:  chunk,
	}

	co := xEnv.Coroutine()
	co.Exdata = code
	co.SetContext(ctx)
	newCodeEv(code, "新增服务").Msg("%s create %s task succeed by %s", way.String(), key, env.Name()).Put()
	return co, code
}

func freeCodeVM(code *Code) {
	ev := newCodeEv(code, "注销服务")
	if e := code.free(); e != nil {
		ev.Msg("失败").E(e).Put()
	} else {
		ev.Msg("成功").Put()
	}
}

func (cd *Code) stop() {
	cd.ToDisable()
}

func (cd *Code) invalid(e error) {
	cd.header.err = e
}

// vela 内部调用
func (cd *Code) vela(name string) *lua.VelaData {
	//判断是否为空
	if cd == nil {
		return nil
	}

	vla, ok := cd.tree.Search(radix.Key(name))
	if ok {
		return vla.(*lua.VelaData)
	}

	return nil
}

func (cd *Code) inMask(name string) bool {
	return in(name, cd.header.mask)
}

func (cd *Code) inLink(name string) bool {
	return in(name, cd.header.link)
}

func (cd *Code) reset() {
	cd.header.mask = cd.header.mask[:0]
	cd.header.link = cd.header.link[:0]
}

func (cd *Code) addMask(name string) {
	if !in(name, cd.header.mask) {
		cd.header.mask = append(cd.header.mask, name)
	}
}

// 记录外连对象
func (cd *Code) addLink(key string) {
	if !in(key, cd.header.link) {
		cd.header.link = append(cd.header.link, key)
	}
}

// 清除code 对象池里面的缓存辅助统计信息
func (cd *Code) clear() error {
	var old []radix.Node

	cd.tree.ForEach(func(node radix.Node) (cont bool) {
		if !cd.inMask(lua.B2S(node.Key())) {
			old = append(old, node)
		}
		return true
	})

	errs := exception.New()
	for _, v := range old {
		cd.tree.Delete(v.Key())
		errs.Try(lua.B2S(v.Key()),
			v.Value().(*lua.VelaData).Close())
	}

	return errs.Wrap()
}

func (cd *Code) free() error {
	cd.cancel()

	errs := exception.New()
	var names []string

	cd.tree.ForEach(func(node radix.Node) (cont bool) {
		key := string(node.Key())
		names = append(names, key)
		ud := node.Value().(*lua.VelaData)
		errs.Try(key, ud.Close())
		return true
	})

	ev := audit.NewEvent("TaskTree.del").From(cd.Key()).Msg("退出服务 %s", strings.Join(names, ","))
	if e := errs.Wrap(); e != nil {
		ev.Subject("退出服务失败").E(e).Alert().Log().Put()
	} else {
		ev.Subject("退出服务成功").E(e).Alert().Log().Put()
	}

	return errs.Wrap()

}

// Close 注销所有
func (cd *Code) Close() error {
	return cd.free()
}

func (cd *Code) foreach(fn func(string, *lua.VelaData) bool) {
	cd.tree.ForEach(func(node radix.Node) (cont bool) {
		key := node.Key()
		val := node.Value().(*lua.VelaData)
		return fn(lua.B2S(key), val)
	})
}

func (cd *Code) Doing() {
	atomic.StoreUint32(&cd.header.status, uint32(vela.Doing))
}

func (cd *Code) ToUpdate() {
	atomic.StoreUint32(&cd.header.status, uint32(vela.Update))
}

func (cd *Code) ToDisable() {
	atomic.StoreUint32(&cd.header.status, uint32(vela.Disable))
}

func (cd *Code) Panic(err error) {
	cd.header.err = err
	atomic.StoreUint32(&cd.header.status, uint32(vela.Panic))
}

func (cd *Code) Fail(err error) {
	cd.reset() //清空
	cd.clear()
	cd.header.err = err
	atomic.StoreUint32(&cd.header.status, uint32(vela.Fail))
}

func (cd *Code) Success() {
	atomic.StoreUint32(&cd.header.status, uint32(vela.Running))
	cd.header.err = nil
	cd.chunk = nil
}

func (cd *Code) T() vela.TaskStatus {
	v := atomic.LoadUint32(&cd.header.status)
	return vela.TaskStatus(v)
}

func (cd *Code) Calling() bool {
	return atomic.AddUint32(&cd.cas, 1) > 1
}

func (cd *Code) pcall(L *lua.LState, way vela.Way) error {
	if cd.Calling() {
		xEnv.Errorf("%s Calling", cd.Key())
		return nil
	}
	defer atomic.StoreUint32(&cd.cas, 0)

	if way != cd.header.way {
		newCodeEv(cd, "唤醒告警").Msg("register: %s but wake up %s", cd.From(), way.String()).Put()
	}

	//切换状态为doing
	cd.Doing()

	fn, err := L.Load(bytes.NewReader(cd.chunk), cd.Key())
	if err != nil {
		cd.Panic(err)
		return cd.Wrap()
	}

	cd.header.uptime = time.Now()
	err = L.CallByParam(cd.header.env.P(fn))
	if err == nil {
		cd.Success()
		return cd.Wrap()
	}

	if ee, ok := err.(*lua.ApiError); ok && strings.HasSuffix(ee.Object.String(), "disable code") {
		cd.ToDisable()
		return cd.Wrap()
	} else {
		cd.Fail(err)
		return cd.Wrap()
	}
}

func (cd *Code) wakeupEv(again bool, way vela.Way) {
	var ev *audit.Event

	if again {
		ev = audit.NewEvent("wakeup.again").From(cd.Key())
	} else {
		ev = audit.NewEvent("wakeup").From(cd.Key())
	}

	er := cd.Wrap()

	if er != nil {
		ev.Subject("唤醒失败")
		ev.Msg("%s wakeup %s task fail", way.String(), cd.Key()).High().E(er).Alert()
	} else {
		ev.Subject("唤醒成功")
		ev.Msg("%s wakeup %s task succeed", way.String(), cd.Key())
	}

	ev.Log().Put()
}

func (cd *Code) wakeup(co *lua.LState, way vela.Way) {

	switch cd.T() {

	case vela.Running:
		return

	case vela.Register:
		_ = cd.pcall(co, way)
		cd.wakeupEv(false, way)

	case vela.Fail:
		_ = cd.pcall(co, way)
		cd.wakeupEv(true, way)

	case vela.Update:
		cd.reset()
		_ = cd.pcall(co, way)
		cd.wakeupEv(true, way)
		_ = cd.clear()

	case vela.Panic:
		audit.NewEvent("wakeup").
			Subject("语法错误").
			Msg("%s wakeup %s fail", way.String(), cd.Key()).
			E(cd.Wrap()).Alert().Log().Put()

	case vela.Doing:
		audit.NewEvent("wakeup").
			Subject("任务加载中").
			Msg("%s wakeup %s fail", way.String(), cd.Key()).
			Alert().Log().Put()

	case vela.Disable:
		cd.reset()
		cd.clear()

	default:
		audit.NewEvent("wakeup").
			Subject("未知状态").
			Msg("%s wakeup %s fail", way.String(), cd.Key()).
			E(fmt.Errorf("task code invalid status , got %s", cd.T().String())).
			Alert().Log().Put()
	}

}

func (cd *Code) Param(key string) (interface{}, bool) {
	if cd.param == nil {
		return nil, false
	}

	v, ok := cd.param[key]
	return v, ok
}
