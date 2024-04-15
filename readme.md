# vela-kit
> vela ssoc security dev kit


## example
```golang
package main

import (
	"github.com/vela-ssoc/vela-kit/vela" //全局抽象接口
	
    kit "github.com/vela-ssoc/vela-kit"  //工具和函数类
    test "github.com/vela-ssoc/vela-test" //自定义模块
)

func injection(xEnv vela.Environment) {
    test.WithEnv(xEnv)
}

func main() {
    deploy := kit.New("vela", kit.All(), kit.Use(injection))
	
	//线上工作模式
	//deploy.Agent()
    
	//调试
    deploy.Debug(kit.Hide{
        Lan:      []string{"ws://1.1.1.1:8082"},
        Hostname: "vela.site.com",
        Edition:  "2.2.0",
        Protect:  true,
    })
	
}
```