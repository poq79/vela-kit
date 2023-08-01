package tasktree

import (
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/vela"
)

type TaskDiffInfo struct {
	Removes []int64      `json:"removes"` // 需要删除的配置名字
	Updates []*TaskEntry `json:"updates"` // 需要执行的配置
}

type TaskEntry struct {
	ID      int64  `json:"id"`
	Dialect bool   `json:"dialect"`
	Name    string `json:"name"`
	Chunk   []byte `json:"chunk"`
	Hash    string `json:"hash"`
}

func (te *TaskEntry) SafeExecute() error {
	if te == nil || te.Name == "" || len(te.Chunk) == 0 {
		return nil
	}
	xEnv.Infof("执行配置: %s hash:%s id:%d", te.Name, te.Hash, te.ID)
	err := root.Tnl(te.ID, te.Name, te.Chunk, xEnv, vela.TRANSPORT, te.Dialect)
	if err != nil {
		xEnv.Errorf("执行配置:%s 失败:%v", te.Name, err)
	}
	return err
}

func (te *TaskEntry) SafeRegister() error {
	return root.Reg(te.ID, te.Name, te.Chunk, xEnv, vela.TRANSPORT, te.Dialect)
}

func (td *TaskDiffInfo) SafeExecute() {
	size := len(td.Updates)
	if size == 0 {
		return
	}

	for i := 0; i < size; i++ {
		entry := td.Updates[i]
		xEnv.Infof("register 配置: %s", entry.Name)
		entry.SafeRegister()
	}

	xEnv.Spawn(0, func() {
		root.Wakeup(vela.TRANSPORT)
	})
}

func (td *TaskDiffInfo) clear() error {
	errs := exception.New()
	for _, id := range td.Removes {
		code := root.FindID(id)
		if code == nil {
			xEnv.Errorf("clear %d task fail error:not found")
			continue
		}
		errs.Try(code.Key(), root.Del(code.Key(), vela.TRANSPORT))
	}
	return errs.Wrap()
}
