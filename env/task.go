package env

import (
	"errors"
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/tasktree"
	"github.com/vela-ssoc/vela-kit/vela"
)

var (
	invalidTaskTreeObject  = errors.New("invalid task tree object")
	notFoundTaskTreeObject = errors.New("not found task tree object")
)

func (env *Environment) WithTaskTree(v interface{}) {
	if env.sub.task != nil {
		return
	}

	tt, ok := v.(*tasktree.TaskTree)
	if ok {
		env.sub.task = tt
		return
	}
	env.Errorf("invalid type task tree object to env , got %T", v)
}

func (env *Environment) TaskSize() int {
	if env.sub.task == nil {
		return 0
	}
	return len(env.sub.task.Keys(""))
}

func (env *Environment) LoadTask(name string, chunk []byte, sess lua.Console) error {
	if env.sub.task == nil {
		return notFoundTaskTreeObject
	}

	return env.sub.task.Load(name, chunk, env, sess)
}

func (env *Environment) DoTask(name string, chunk []byte, way vela.Way) error {
	if env.sub.task == nil {
		return notFoundTaskTreeObject
	}

	return env.sub.task.Do(name, chunk, env, way)
}

func (env *Environment) DoTaskFile(path string, way vela.Way) error {
	if env.sub.task == nil {
		return notFoundTaskTreeObject
	}

	return env.sub.task.DoFile(path, env, way)
}

func (env *Environment) RegisterTask(id int64, name string, chunk []byte, way vela.Way, dialect bool) error {
	if env.sub.task == nil {
		return notFoundTaskTreeObject
	}

	return env.sub.task.Reg(id, name, chunk, env, way, dialect)
}

func (env *Environment) WakeupTask(way vela.Way) error {
	if env.sub.task == nil {
		return notFoundTaskTreeObject
	}
	return env.sub.task.Wakeup(way)
}

func (env *Environment) RemoveTask(name string, way vela.Way) error {
	if env.sub.task == nil {
		return notFoundTaskTreeObject
	}

	return env.sub.task.Del(name, way)
}

func (env *Environment) DoTaskByTnl(id int64, key string, chunk []byte, way vela.Way, dialect bool) error {
	if env.sub.task == nil {
		return notFoundTaskTreeObject
	}

	return env.sub.task.Tnl(id, key, chunk, env, way, dialect)
}

func (env *Environment) TaskList() []*vela.Task {
	if env.sub.task == nil {
		return nil
	}

	return env.sub.task.ToTask()
}

func (env *Environment) FindCode(name string) vela.Code {
	if env.sub.task == nil {
		return nil
	}

	return env.sub.task.Code(name)
}

func (env *Environment) FindTask(name string) *vela.Task {
	tl := env.TaskList()
	if len(tl) == 0 {
		return nil
	}

	for _, t := range tl {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (env *Environment) FindProc(cname string, name string) (*lua.VelaData, error) {
	if env.sub.task == nil {
		return nil, fmt.Errorf("not found valid task in Environment")
	}

	code := env.FindCode(cname)
	if code == nil {
		return nil, fmt.Errorf("not found %s code", code)
	}

	proc := code.Get(name)
	if proc == nil {
		return nil, fmt.Errorf("not found %s.%s", code, name)
	}

	return proc, nil
}
