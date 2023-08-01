# audit
全局事件审计模块,主要用作重要事件处理需求


## vela.event
- ev = vela.event(typeof)
- ev = vela.event{typeof , msg , remote , port ,...}
- 采用的是链式调用的方式和table初始化方法
### 字段
- 满足index 和new index 接口
- [time]()
- [id]()
- [inet]()
- [subject]()
- [addr]()
- [port]()
- [from]()
- [typeof]()
- [user]()
- [auth]()
- [msg]()
- [err]()
- [region]()
- [alert]()
- [level]()

### 函数接口 函数支持链式调用
- [Time(v)]()   时间
- [Subject(v)]()主题
- [Remote(v)]() 远程地址
- [Port(v)]()   远程端口
- [From(v)]()   来源
- [Typeof(v)]() 类型
- [User(v)]() 用户
- [Auth(v)]() 认证信息
- [Msg(v)]() 事件消息
- [E(v)]()  报错
- [Region(v)]() 地理位置
- [Alert(v)]()  是否告警
- [Level(n)]()  等级
- [Log()]()     打印日志
- [Put(b , b , n)]() 是否提交 参数1: 是打印记录日志  参数2： 是否告警  参数3： 设置等级

```lua
    local ev = vela.event("demo").Msg("helo").Port(1)
        .Remote("127.0.0.1").Auth("use:a pass:2").E("fail")
    ev.Put(true , true)
```

## 注意
默认如果 alert ~= true 系统就会发生告警