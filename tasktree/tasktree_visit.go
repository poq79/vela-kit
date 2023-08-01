package tasktree

import (
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"io/ioutil"
)

type RequestTask struct {
	Tasks []*vela.Task `json:"tasks"`
}

func (tt *TaskTree) Name() string {
	return "TaskTree"
}

func (tt *TaskTree) Wakeup(way vela.Way) error {
	ev := audit.NewEvent("TaskTree.wakeup").From(way.String())
	err := tt.wakeup(way)
	if err != nil {
		ev.Subject("唤醒服务失败").E(err).High().Alert().Put()
		return err
	}
	ev.Subject("唤醒服务成功").Notice().Alert().Put()
	return nil
}

func (tt *TaskTree) Reg(id int64, name string, chunk []byte, env vela.Environment, way vela.Way, dialect bool) error {
	return wrap(root.reg(id, name, chunk, env, way, dialect))
}

func (tt *TaskTree) Code(name string) vela.Code {
	_, code := root.GetCodeVM(name)
	return code
}

func (tt *TaskTree) Vela(cname string, proc string) (*lua.VelaData, error) {
	code := tt.code(cname)
	if code == nil {
		return nil, NotFoundE
	}

	obj := code.vela(proc)
	if obj == nil {
		return nil, NotFoundE
	}

	return obj, nil
}

func (tt *TaskTree) Del(name string, way vela.Way) error {

	//Del 删除需要关闭的code

	ev := audit.NewEvent("TaskTree.del").
		Msg("删除: %s 服务 来源: %s", name, way.String()).
		From(name)

	if e := root.remove(name); e != nil {
		ev.Subject("删除服务失败").E(e).Log().Alert().High().Put()
		return e
	} else {
		ev.Subject("删除服务成功").Log().Alert().High().Put()
		return nil
	}
}

// Do 运行单个代码块
func (tt *TaskTree) Do(key string, chunk []byte, env vela.Environment, way vela.Way) error {
	co := tt.reg(0, key, chunk, env, way, true)
	return wakeup(co, way)
}

func (tt *TaskTree) Tnl(id int64, key string, chunk []byte, env vela.Environment, way vela.Way, dialect bool) error {
	co := tt.reg(id, key, chunk, env, way, dialect)
	return wakeup(co, way)
}

// DoFile 外部加载代码
func (tt *TaskTree) DoFile(path string, env vela.Environment, way vela.Way) error {
	chunk, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return tt.Do(lua.FileSuffix(path), chunk, env, way)
}
func (tt *TaskTree) Marshal() []byte {
	return root.Byte()
}

func (tt *TaskTree) Load(name string, chunk []byte, env vela.Environment, sess lua.Console) error {
	ev := audit.NewEvent("TaskTree.load").From(name).Msg("来源: vela-console")
	err := root.load(name, chunk, env, sess)
	if err != nil {
		ev.Subject("加载服务失败").E(err).High().Alert().Put()
		return err
	}

	ev.Subject("加载服务成功").E(err).High().Alert().Put()
	return nil
}

func (tt *TaskTree) ToTask() []*vela.Task {
	var ats []*vela.Task

	tt.forEach(func(key string, co *lua.LState, code *Code) bool {
		at := &vela.Task{
			Name:    code.Key(),
			Link:    code.Link(),
			Hash:    code.Hash(),
			From:    code.From(),
			Status:  code.Status(),
			Uptime:  code.header.uptime,
			ID:      code.header.id,
			Dialect: code.header.dialect,
		}
		at.Runners = code.List()

		if code.header.err == nil {
			at.Failed = false
		} else {
			at.Cause = code.header.err.Error()
		}
		ats = append(ats, at)
		return true
	})
	return ats
}
