package env

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type Startup struct {
	Node    Node      `json:"node"`
	Logger  Logger    `json:"logger"`
	Console Console   `json:"console"`
	Extends []*Extend `json:"extends"`
}
type Node struct {
	DNS    string `json:"dns"`
	Prefix string `json:"prefix"`
}

type Logger struct {
	Level    string `json:"level"` // 日志级别 debug/info/error
	Filename string `json:"filename"`
	Console  bool   `json:"console"`
	Format   string `json:"format"` // 日志格式 text/json
	Caller   bool   `json:"caller"` // 是否打印调用函数名字
	Skip     int    `json:"skip"`
}

type Console struct {
	Enable  bool   `json:"enable"`
	Network string `json:"network"`
	Address string `json:"address"`
	Script  string `json:"script"`
}

type Extend struct {
	Name  string `json:"name"`
	Type  string `json:"type"` // number bool string ref string_readonly
	Value string `json:"value"`
}

func (env *Environment) StartupHandler(ctx *fasthttp.RequestCtx) error {
	body := ctx.Request.Body()
	if len(body) == 0 {
		return fmt.Errorf("startup config got nil")
	}

	var startup Startup
	err := json.Unmarshal(body, &startup)
	if err != nil {
		return err
	}

	r := env.R()
	_, err = r.Call("/api/v1/inline/agent/logger", startup.Logger)
	if err != nil {
		env.Errorf("call logger startup fail %v", err)
	}
	_, err = r.Call("/api/v1/inline/agent/node", startup.Node)
	if err != nil {
		env.Errorf("call node startup fail %v", err)
	}

	_, err = r.Call("/api/v1/inline/agent/extends", startup.Extends)
	if err != nil {
		env.Errorf("call extends startup fail %v", err)
	}

	return err
}

func (env *Environment) startup() {
	r := env.R()
	r.POST("/api/v1/agent/startup", env.Then(env.StartupHandler))
}
