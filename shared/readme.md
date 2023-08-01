# shared
添加shared arc 缓存算法 用来临时存储数kv 结构数据 是何小量数据

## rock.shared.* = size 
新增缓存空间 大小为key的个数，而非内存大小
```lua
    rock.shared.shm = 300 --名称为shm
```

## rock.shared.*.set(key , val , [ttl]) 
- 存储一个新的key 如果key存在的话 ，会直接覆盖， 设置了ttl的话，会重新计算时间， 没有沿用上次的expire
- key: string
- val: 非nil对象
- ttl: 单位 millisecond(微秒) , 默认为永久
```lua

    local shm = rock.shared.shm
    shm.set("192.168.1.1" , 2 , 4000) -- 4s
    
    shm.set("192.168.1.1" , 30) -- 2个参数时候
    
```

## rock.shared.*.incr(key , step , [ttl])
- 加法运算 如果key不存在 默认创建， 此时默认数值为step , ttl 跟set 方法类似 
```lua

    local shm = rock.shared.shm
    shm.incr("192.168.1.1" , 1 , 502)
```

## rock.shared.*.del(key)
- 删除key
```lua
    local shm = rock.shared.shm
    shm.del("192.168.1.1" , 1)
```

## rock.shared.*.count
- 当前缓存中的数据条数， 不包括已经超时的
```lua
    rock.ERR(rock.shared.shm.count)
```

## rock.shared.*.count_all
- 当前缓存中的数据条数， 包括已经超时的
```lua
    rock.ERR(rock.shared.shm.count_all)
```

## rock.shared.pairs(fn)
- 迭代遍历
```lua
    rock.shared.shm.pairs(function(item , stop)
        rock.ERR("key: " , item.key)
        rock.ERR("val: " , item.val)
        rock.ERR("clock: " , item.clock)
        rock.ERR("expire:" , item.expire)
    end)
```
