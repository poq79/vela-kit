package tasktree

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

type WebConsole struct {
	ctx *fasthttp.RequestCtx
}

func (w *WebConsole) Println(v string) {
	w.ctx.Write(lua.S2B(v))
	w.ctx.WriteString("\n")
	return
}

func (w *WebConsole) Printf(f string, v ...interface{}) {
	chunk := fmt.Sprintf(f, v...)
	w.ctx.Write(lua.S2B(chunk))
	return
}

func (w *WebConsole) Invalid(f string, v ...interface{}) {
	chunk := fmt.Sprintf(f, v...)
	w.ctx.Write(lua.S2B(chunk))
	return
}

func (tt *TaskTree) DefineViewTask(ctx *fasthttp.RequestCtx) error {
	chunk, err := json.Marshal(RequestTask{Tasks: tt.ToTask()})
	if err != nil {
		//r.Bad(ctx, http.StatusInternalServerError, problem.Title("内部错误"), problem.Detail("序列化失败:%v", err))
		return err
	}

	ctx.Write(chunk)
	return nil
}

func (tt *TaskTree) diff(ctx *fasthttp.RequestCtx) error {

	body := ctx.Request.Body()

	if len(body) == 0 {
		return tt.DefineViewTask(ctx)
	}

	var diff TaskDiffInfo
	err := json.Unmarshal(body, &diff)
	if err != nil {
		return err
	}

	diff.SafeExecute()
	if e := diff.clear(); e != nil {
		xEnv.Errorf("task clear fail %v", e)
	}

	return tt.DefineViewTask(ctx)
}

func (tt *TaskTree) cli(ctx *fasthttp.RequestCtx) error {
	name := ctx.QueryArgs().Peek("name")
	body := ctx.Request.Body()
	if len(body) == 0 {
		return fmt.Errorf("not found body")
	}

	if len(name) == 0 {
		return fmt.Errorf("not found name")
	}

	err := tt.Load(string(name), body, xEnv, &WebConsole{ctx: ctx})
	if err != nil {
		return err
	}

	return nil
}

func (tt *TaskTree) delete(ctx *fasthttp.RequestCtx) error {
	name := ctx.QueryArgs().Peek("name")
	if len(name) == 0 {
		return fmt.Errorf("not found task name %s", string(name))
	}

	c := tt.code(string(name))
	if c == nil {
		return fmt.Errorf("not found task %s", string(name))
	}

	if c.header.way != vela.CONSOLE {
		return fmt.Errorf("not allow delete")
	}

	tt.del(string(name))
	tt.Report()
	return nil
}

/*
	# param
	{
		"name": "vela-abc",
		"code": "print(hello)",
		"param": {
			arr: {"name" , "name2" , "name3"},
			cnt: 10,
			ip: "10.205.14.127",
		}
		"report": true
	}
*/

func (tt *TaskTree) scanner(ctx *fasthttp.RequestCtx) error {
	name := ctx.QueryArgs().Peek("name")
	body := ctx.Request.Body()
	if len(body) == 0 {
		return fmt.Errorf("not found body")
	}

	if len(name) == 0 {
		return fmt.Errorf("not found name")
	}

	data := make(map[string]interface{}, 32)
	err := tt.Scan(xEnv, 0, string(name), body, data, 30)
	if err != nil {
		return err
	}

	return nil
}

func (tt *TaskTree) console(ctx *fasthttp.RequestCtx) error {
	body := ctx.Request.Body()
	if len(body) == 0 {
		return fmt.Errorf("not found body")
	}

	var tx Console
	if err := json.Unmarshal(body, &tx); err != nil {
		return err
	}

	err := tt.Load("console", body, xEnv, &WebConsole{ctx: ctx})
	if err != nil {
		return err
	}
	return nil
}

func (tt *TaskTree) define(r vela.Router) {
	_ = r.POST("/api/v1/agent/scan/run", xEnv.Then(tt.scanner))
	_ = r.POST("/api/v1/agent/task/diff", xEnv.Then(tt.diff))
	_ = r.POST("/api/v1/agent/task/status", xEnv.Then(tt.DefineViewTask))
	_ = r.POST("/api/v1/arr/agent/task/status", xEnv.Then(tt.DefineViewTask))
	_ = r.POST("/api/v1/arr/agent/task/cli", xEnv.Then(tt.cli))
	_ = r.DELETE("/api/v1/arr/agent/task/cli", xEnv.Then(tt.delete))
	_ = r.POST("/api/v1/arr/agent/task/console", xEnv.Then(tt.console))
}
