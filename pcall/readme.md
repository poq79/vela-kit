# pcall
rock-go safe调用接口防止 系统崩溃 只提供go API接口

## 基础使用
- [pcall.Exec(func , ...interface) *safe](#)
- [safe.Ok(func)](#)
- [safe.Err(func)](#)
- [safe.Time(v)](#)
- [safe.Spawn(v)](#)

```go
    package demo

import (
	"fmt"
	"github.com/vela-ssoc/pivot/pcall"
	"time"
)

func a() {
	fn := func(s string) {
		fmt.Errorf("%s", s)
	}

	pcall.Exec(fn).                             //注入方法
		Ok(func() { print("ok") }).             //执行完成 后动作
		Err(func(err error) { print("fail") }). //报错执行 动作
		Time(5 * time.Second).                  //执行timeout时间
		Spawn()                                 //开始执行

}

```