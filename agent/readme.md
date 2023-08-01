# rock
磐石系统框架代码

## 内置库说明
### rock.console
启用console调试功能
功能说明:[rock.console](https://github.com/vela-ssoc/pivot/tree/master/console)

### rock.node
设置基础参数
功能说明:[rock.node](https://github.com/vela-ssoc/pivot/tree/master/node)

### rock.shared
设置基础参数
功能说明:[rock.shared](https://github.com/vela-ssoc/pivot/tree/master/shared)

### rock.load
设置基础参数
功能说明:[rock.load](https://github.com/vela-ssoc/pivot/tree/master/plugin)

### rock.event
默认时间审计
功能说明:[rock.event](https://github.com/vela-ssoc/pivot/tree/master/audit)

### rock.require
导入三方依赖
功能说明:[rock.require](https://github.com/vela-ssoc/pivot/tree/master/require)

### rock.pipe
pipe内部处理接口
功能说明:[rock.require](https://github.com/vela-ssoc/pivot/tree/master/pipe)

### bbolt.bucket
内部key-val地址库
功能说明:[bbolt.bucket](https://github.com/vela-ssoc/pivot/tree/master/bucket)

### rock.hm
内部hashmap
功能说明:[rock.hm](https://github.com/vela-ssoc/pivot/tree/master/hashmap)

### service
保存所有PROC服务树 , 采用radix基数树存储 协程安全
功能说明:[rock.service](https://github.com/vela-ssoc/pivot/tree/master/service)

### rock.xreflect
一个完全映射的接口 利用反射可以完全暴露golang中的对象
功能说明: [rock.xreflect](https://github.com/vela-ssoc/pivot/tree/master/xreflect)

### rock.file
用来新建一个file的文件追加操作 一般是日志输出的位置
功能说明: [rock.file](https://github.com/vela-ssoc/pivot/tree/master/file)

### rock.json
通用的JSON的解码和编码
功能说明: [rock.json](https://github.com/vela-ssoc/pivot/tree/master/json)

### rock.logger
设置框架日志参数
功能说明: [rock.logger](https://github.com/vela-ssoc/pivot/tree/master/logger)

### rock.region
获取IP地址位置信息
功能说明: [rock.region](https://github.com/vela-ssoc/pivot/tree/master/region)

### http
内置http请求库
功能说明: [http](https://github.com/vela-ssoc/pivot/tree/master/request)

### std
内置std控制台函数,全局变量
功能说明: [std](https://github.com/vela-ssoc/pivot/tree/master/xlib#std)

### timer
内置时间模块
功能说明: [timer](https://github.com/vela-ssoc/pivot/tree/master/xlib#timer)

# go module 私有库配置
让go get 可以直接下载私有库

### 1.设置GOPRIVATE
```shell
    go env -w GOPRIVATE=github.com/vela-ssoc
```

### 2.在 ~/.gitconfig 中添加
```shell
[url "git@github.com:"]
	insteadOf = https://github.com/
```

### 3.如果配置了GOPROXY不要加入以下配置
```plain
[url "ssh://git@github.com/"]
        insteadOf = https://github.com/
```
### 4.在 ~/.netrc 中添加(非必须)
```plain
machine github.com login <GitHub的用户名> password <GitHub的KEY>
```

### 5.配置Public Key
```shell
# 生成公钥
ssh-keygen -t rsa -C "email@xxx.com"
# 将生成的 id_rsa.pub 配置到 GitHub 中
```
### 6.测试
```shell
ssh -T git@github.com
```

# 结构说明
```go
package main

import (
    "github.com/vela-ssoc/rock"
    "github.com/vela-ssoc/pivot/xcall"
    
    //自定义的模块
    cron "github.com/vela-ssoc/rock-cron-go"
)

func init() {
	rock.Inject(xcall.Rock ,  cron.Constructor )
	//可以添加自定义的模块
	//todo
}

func main() {
	//启动函数
	rock.Setup(xcall.Rock)
}
```

## 复杂结构
主要是用来对接 struct 和 lua虚拟机方法绑定

### userdata
官方内置方法，可以查看官方操作手册

### lightuserdata
内置的proc绑定方法，一般实现固定接口，主要用于服务对象
这个时候A必须要满足lua.LightUserDataIFace
```golang
    type A struct {
        lua.Super	
        name string
    }
    
    a := &A{name: "123"}
    ud := lua.NewLightUserData(a)
```
### context
lua中内置的context的处理函数
```lua
    local ctx , cancel = context.WithTimeout(5) --单位是秒
    local ctx = context.Background()

```
### AnyData
内置的万能数据类型 减少内存之间的交互需求， 用法如下
这个方法可以绕过对象的公私有对象
```go
    //自动查找反射函数
    type A struct {
        name string `lua:"name"`
        pass string `lua:"pass"`
    }
    
    func (a *A) Name(L *lua.LState) lua.LValue {
    	return lua.LString(a.name)
    }
    
    func (a *A) Pass(L *lua.LState) int {
    	L.Push(lua.LString(a.pass))
    	return 1
    }
    
    func new(L *lua.LState) int {
    	L.Push(L.NewAnyData(new(A)))
        return 1	
    }

    //拦截器会自动查找 对应的方法 如： func(*Lstate) int 或者 func(*LState) LValue
    /* 
        local a = new(A)
        a.Name()
        a.Pass()
     */
```
```go
    //通过用户定义的方法查找 优化性能
    type A struct {
    	lua.NoReflect
    	name string
    }
   
    //需要自定义Get方法
    func (a *A) Get( key string ) lua.LValue {
    	if key == name {
    	    return lua.LString(a.name)	
        }
    }
    
    func (a *A) Set( L *lua.LState , key string , val lua.LValue) {
    	
    }
    
    //lua中的用法同上

```
)
