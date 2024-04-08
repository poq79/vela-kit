package tasktree

import (
	"errors"
	"fmt"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/pcall"
	"github.com/vela-ssoc/vela-kit/vela"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vela-ssoc/vela-kit/lua"
)

var (
	NotFoundE = errors.New("not found task")
	root      = &TaskTree{scans: new(sync.Map), tasks: new(sync.Map)}
)

type TaskTree struct {
	tasks  *sync.Map
	scans  *sync.Map
	waking uint32
}

func (tt *TaskTree) String() string                         { return lua.B2S(tt.Byte()) }
func (tt *TaskTree) Type() lua.LValueType                   { return lua.LTObject }
func (tt *TaskTree) AssertFloat64() (float64, bool)         { return 0, false }
func (tt *TaskTree) AssertString() (string, bool)           { return lua.B2S(tt.Byte()), true }
func (tt *TaskTree) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (tt *TaskTree) Peek() lua.LValue                       { return tt }

func (tt *TaskTree) forEach(callback func(key string, co *lua.LState, code *Code) bool) {
	tt.tasks.Range(func(key, v interface{}) bool {
		co := v.(*lua.LState)
		code, ok := CheckCodeVM(co)
		if !ok {
			tt.tasks.Delete(key)
			return true
		}
		return callback(key.(string), co, code)
	})
}

func (tt *TaskTree) FindID(id int64) (ret *Code) {
	tt.forEach(func(key string, co *lua.LState, code *Code) bool {
		if code.header.id == id {
			ret = code
			return false // break foreach
		}
		return true // true: next
	})
	return
}

func (tt *TaskTree) BeLinked(name string) (links []string) {
	tt.forEach(func(key string, co *lua.LState, code *Code) bool {
		if code.inLink(name) {
			links = append(links, code.Key())
		}
		return true
	})
	return
}

func (tt *TaskTree) remove(name string) error {

	links := tt.BeLinked(name)
	if len(links) > 0 {
		return errors.New("find " + name + " code link in " + strings.Join(links, ","))
	}
	tt.del(name)
	return nil
}

func (tt *TaskTree) Close() error {
	tt.forEach(func(key string, co *lua.LState, code *Code) bool {
		pcall.Exec(code.Close).Time(time.Second * 30).Spawn()
		return true
	})

	return nil
}

func (tt *TaskTree) Byte() []byte {
	buf := kind.NewJsonEncoder()
	buf.Arr("")
	tt.forEach(func(key string, co *lua.LState, cd *Code) bool {
		buf.Tab("")
		buf.KV("key", cd.Key())
		buf.KV("link", cd.Link())
		buf.KV("status", cd.Status())

		buf.KV("hash", cd.Hash())
		buf.KV("way", cd.From())
		buf.KV("uptime", cd.Time())

		if e := cd.Wrap(); e != nil {
			buf.KV("err", e.Error())
		} else {
			buf.KV("err", nil)
		}

		buf.Arr("vela")
		cd.foreach(func(name string, ud *lua.VelaData) bool {
			buf.Tab("")
			buf.KV("name", name)
			buf.KV("type", ud.Data.Type())
			buf.KV("state", ud.Data.State().String())
			buf.End("},")
			return true
		})
		buf.End("]},")
		return true
	})
	buf.End("]")
	return buf.Bytes()
}

func (tt *TaskTree) GetCodeVM(name string) (*lua.LState, *Code) {
	v, ok := tt.tasks.Load(name)
	if !ok {
		return nil, nil
	}

	co := v.(*lua.LState)

	code, ok := CheckCodeVM(co)
	if ok {
		return co, code
	}

	tt.tasks.Delete(name)
	return nil, nil
}

func (tt *TaskTree) code(name string) *Code {
	_, code := tt.GetCodeVM(name)
	return code
}

func (tt *TaskTree) insert(name string, co *lua.LState) {
	tt.tasks.Store(name, co)
}

func (tt *TaskTree) Keys(prefix string) []string {

	var keys []string

	tt.forEach(func(name string, co *lua.LState, code *Code) bool {
		if prefix == "" {
			keys = append(keys, name)
			return true
		}

		if strings.HasPrefix(name, prefix) {
			keys = append(keys, name)
		}
		return true
	})
	return keys
}

func (tt *TaskTree) del(name string) {
	co, code := tt.GetCodeVM(name)
	if co == nil || code == nil {
		xEnv.Errorf("del task %s fail %v", name, NotFoundE)
		return
	}
	defer code.Close()
	tt.tasks.Delete(name)
	//freeCodeVM(code)
}

func (tt *TaskTree) reg(id int64, cname string, chunk []byte, env vela.Environment, way vela.Way, dialect bool) *lua.LState {
	co, code := tt.GetCodeVM(cname)
	if co == nil {
		co, code = newCodeVM(cname, chunk, env, way)
		tt.insert(cname, co)
		newCodeEv(code, "添加服务").Msg("%s 注册成功", cname)
		goto done
	}

	code.ToUpdate()
	code.Update(co, chunk, env, way)
	newCodeEv(code, "更新服务").Msg("%s 注册成功", cname)

done:
	code.header.id = id
	code.header.dialect = dialect

	return co
}

func (tt *TaskTree) Report() {

	task := RequestTask{tt.ToTask()}
	if e := xEnv.Push("/api/v1/broker/task/status", task); e != nil {
		xEnv.Errorf("report task call fail %v", e)
	} else {
		//xEnv.Infof("report task call succeed %v", task)
	}
}

func (tt *TaskTree) wakeup(way vela.Way) error {
	if tt.Waking() {
		return fmt.Errorf("ssc task tree waking up")
	}

	defer tt.deferFn()

	errs := exception.New()
	tt.forEach(func(key string, co *lua.LState, code *Code) bool {
		err := wakeup(co, way)
		if err != nil {
			xEnv.Errorf("wakeup 配置：%s 失败:%v", key, err)
			errs.Try(key, err)
		}
		return true
	})

	tt.Report()
	return errs.Wrap()
}

// load console代码
func (tt *TaskTree) load(key string, chunk []byte, env vela.Environment, sess lua.Console) error {
	co := tt.reg(0, key, chunk, env, vela.CONSOLE, true)
	co.Console = sess
	defer func() {
		co.Console = nil
	}()
	return wakeup(co, vela.CONSOLE)
}

// ExistProc 判断服务对象是否存在
func (tt *TaskTree) ExistProc(proc string) bool {
	ret := false
	tt.forEach(func(key string, co *lua.LState, code *Code) bool {
		if code.Exist(proc) {
			ret = true
			return false //终止 foreach
		}
		return true
	})
	return ret
}

// ExistCode 判断服务代码快是否存在
func (tt *TaskTree) ExistCode(name string) bool {
	_, code := tt.GetCodeVM(name)
	return code != nil
}

func (tt *TaskTree) Waking() bool {
	return atomic.AddUint32(&tt.waking, 1) > 1
}

func (tt *TaskTree) deferFn() {
	atomic.StoreUint32(&tt.waking, 0)
}

/*
func (tt *TaskTree) again() {
	if tt.Waking() {
		audit.NewEvent("task.again").Subject("定时任务").Msg("同步服务状态").Put()
		return
	}

	defer tt.deferFn()

	errs := exception.New()
	tt.forEach(func(key string, co *lua.LState, code *Code) bool {
		errs.Try(key, wakeup(co, vela.AGAIN))
		return true
	})

	if e := errs.Wrap(); e != nil {
		xEnv.Errorf("task wakeup again fail %v", e)
	}
}
*/
