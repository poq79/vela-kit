package third

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/vela"
	"os"
)

func (th *third) OnChange(c Change) {
	if c.Name == "" {
		return
	}

	switch c.Event {
	case "update":
		info, ok := th.info(c.Name)
		if !ok {
			return
		}
		th.update(c.Name, info.Hash)
	case "delete":
		th.drop(&vela.ThirdInfo{Name: c.Name})

	default:
		xEnv.Errorf("invalid third event %v", c)
	}
}

func (th *third) check() {
	err := os.Mkdir(th.dir, 0755)
	if err != nil {
		xEnv.Errorf("create 3rd fail %v", err)
	}
}

func (th *third) httpHandleInfo(ctx *fasthttp.RequestCtx) error {
	name, ok := ctx.UserValue("name").(string)
	if !ok || name == "" {
		return fmt.Errorf("not found param name")
	}

	info, ok := th.info(name)
	if ok {
		ctx.Write(info.Byte())
		return nil
	}

	return fmt.Errorf("not found %s", name)
}

func (th *third) httpHandleLoad(ctx *fasthttp.RequestCtx) error {
	name, ok := ctx.UserValue("name").(string)
	if !ok || name == "" {
		return fmt.Errorf("not found param name")
	}

	info, err := th.load(name)
	if err != nil {
		return err
	}

	ctx.Write(info.Byte())
	return nil
}

func (th *third) httpHandleClear(ctx *fasthttp.RequestCtx) error {
	th.clear()
	return nil
}

func (th *third) httpHandleDiff(ctx *fasthttp.RequestCtx) error {
	body := ctx.Request.Body()
	var c Change
	if len(body) == 0 {
		return fmt.Errorf("not found third info")
	}
	err := json.Unmarshal(body, &c)
	if err != nil {
		return err
	}
	go th.OnChange(c)
	ctx.WriteString("ok")
	return nil
}

func (th *third) define() {
	r := xEnv.R()
	r.GET("/inline/third/info/{name}", xEnv.Then(th.httpHandleInfo))
	r.GET("/inline/third/load/{name}", xEnv.Then(th.httpHandleLoad))
	r.GET("/inline/third/clean", xEnv.Then(th.httpHandleClear))
	r.POST("/api/v1/agent/third/diff", xEnv.Then(th.httpHandleDiff))
}
