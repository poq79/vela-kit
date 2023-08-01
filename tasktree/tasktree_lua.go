package tasktree

import (
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

// Index 获取服务对应的代码块
func (tt *TaskTree) Index(L *lua.LState, cname string) lua.LValue {
	cd, ok := CheckCodeVM(L)
	if !ok {
		goto done
	}

	//自己调用自己
	if cd.Key() == cname {
		audit.NewEvent("task").Subject("循环引用").
			Msg("%s loop call %s", L.CodeVM(), cname).Put()

		L.RaiseError("loop call %s", cname)
		return nil

	}

	cd.addLink(cname)

done:
	co, code := tt.GetCodeVM(cname)
	if co == nil {
		L.RaiseError("not found %s", cname)
		return lua.LNil
	}

	if code.IsReg() {
		wakeup(co, vela.INLINE)
	}

	return code
}

func servLuaInjectApi(env vela.Environment) {
	env.Global("task", root)
}
