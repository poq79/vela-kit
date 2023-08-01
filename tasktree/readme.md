# service
用来存储所有rock-go启动的服务信息,线程安全的radix基数数存储 

## 外部调用
- 主要是在lua中怎么调用其他CODE中的服务
- 注意命名虽然不做要求 最好用特殊符号下划线分割
```lua
    --获取启动的服务代码块
    local cd = task.fasthttp 
    local cd = task["fasthttp"]
    
    --获取代码中的服务
    local proc = cd.kfk
    local proc = cd["kfk"]

    --其次启动
    proc.start(kfk , ud2 , ud3)

    --只允许本地调用 不允许导出
    proc.inline(kfk)
```