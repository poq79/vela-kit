# json
磐石系统JSON序列化

## rock.json.encode
通用的JSON编码
```lua
    local encode = rock.json.encode
    local v = encode({['name']="x" , ['age'] = 1 , ['s']={3,4,5}})

    --local ud = userdata
    local v = encode(ud)
    
    --local ud = lightuserdata
    local v = encode(ud)
```
## rock.json.decode
通用的JSON解码
```lua
    local decode = rock.json.decode
    local t = decode('{"name":"x" , "pass":"y"}')
```

## rock.fastJson(v)
fastjson解码
```lua
    local f , e = rock.fastJson('{"name":"x" , "pass":"Y" , "age": 18 , "up": true}')
    if e ~= nil then
        --todo     
    end
    
    local name = f.str("name")
    local pass = f.str("pass")
    local age = f.int("age")
    local up = f.bool("up")
```