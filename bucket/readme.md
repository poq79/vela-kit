# vela-db
db 本地key-valule 数据存储接口 利用的bbolt 的本地存储数据库

## 基础用法
基础常用接口说明

### db.bucket(string , string , ...)
打开一个bucket 多个代表桶的级联
```lua
    local bkt = bbolt.bucket("2022-02-21" , "record")
```

### bucket基础用法

- [bucket.byte](bucket.byte)
- [bucket.export](bucket.export)
- [bucket.get](bucket.get)
- [bucket.set](bucket.set)
- [bucket.delete](bucket.delete)
- [bucket.remove](bucket.remove)
- [bucket.pairs](bucket.pairs)
- [bucket.info](bucket.info)
- [bucket.count](bucket.count)
- [bucket.depth](bucket.depth)
- [bucket.incr](bucket.incr)
- [bucket.suffix](bucket.suffix)
- [bucket.prefix](bucket.prefix)
```lua
    local bucket = bbolt.bucket("2022-02-21" , "record")
    bucket.byte()
    bucket:get("key")
    bucket:set("key" , "helo")
    --todo
```

#### bucket.byte
- bbolt bucket 数据导出接口 用法
```lua
    local bucket = bbolt.bucket("xxx")
    local data = bucket.byte()
```

#### bucket.export
- bbolt expore 数据导出格式校验
```lua
    local bucket = bbolt.bucket("xxx")
    bucket.export("json")
    bucket.byte() --导出接送结构
```

#### bucket.get
- 获取bucket中的数据
```lua
    local bucket = bbolt.bucket("xx")
    local value , err = bucket:get("key")
```

#### bucket.set
- 修改bucket中的数据 
- err = bucket:set(key , value , ttl)
```lua
    local err = bucket:set("helo" , 10 , 1000)
```

#### bucket.delete
- 删除bucket中值
- err = bucket:delete(key , key , ...)
```lua
    local err = bucket:delete("key1" , "key2")
```
#### bucket.remove
- 删除bucket中的桶
- err = bucket:remove(bkt , bkt, ...)
```lua
    local err = bucket:remove("bkt1" , "bkt2")
```

#### bucket.paris
- 遍历bucket的数据
- err = bucket:paris(function(key , val) end)
```lua
    local err = bucket:paris(function(key , val) 
        print(key)
        print(val)
    end)
```

#### bucket.info
- 获取bucket的状态
- info = bucket:info()
```lua
    local info = bucket:info()
```

#### bucket.count
- 获取bucket中key的总数
- num = bucket:count()
```lua
    local num = bucket:count()
```

#### bucket.depth
- 获取bucket的深度
- d = bucket:depth()
```lua
    local d = bucket:depth()
```

#### bucket.incr
- 将bucket中key的值做加法 并返回
- v = bucket:incr(key , val , init , ttl)
```lua
    bucket:set("_int" , 0)
    bucket:incr("_int" , 1)
    
    bucket:incr("aaa" , 1 , 10 , 100) -- 100*millisecond expire
```

#### bucket.suffix
- 遍历bucket中suffix中后缀的值
- tab = bucket:suffix(string)
```lua
    local tab = bucket:suffix("x-")

    for k , v in pairs(tab) do
        print(k)
        print(v)
    end 
```

#### bucket.prefix
- 遍历bucket中prefix中后缀
- 用法同上

### bucket默认接口
- 默认情况下满足index 和newindex 的接口
```lua
    local bucket = bbolt.bucket("2022-02-21" , "k")
    bucket._string = "helo"   --key: _string
    bucket._int = 10          --key: _int
    bucket._bool = false      --key: false
    
    print(bucket._string)
    print(bucket._int)
    print(bucket._bool)
```

## 其他语法
- 其他接口

### bbolt.info
- 获取一级bucket的名称
- name = bbolt.info()
```lua
    local names = bbolt.info()
```

### bbolt.remove
- 删除一级bucket
- bbolt.remove(name)
```lua
    bbolt.remove("xxx")
```