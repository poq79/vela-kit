# node 
存储服务器信息

## 函数
```lua
    rock.node{
        id = "x-xx-x-0", --设置服务器ID
        resolve = "114.114.114.114:53", --设置dns解析地址
    }

    local id = rock.ID() -- 获取ID
    local inet = rock.inet() -- 获取IP地址
```